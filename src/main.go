package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"src/api"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAdd string
}

func main() {
	port := os.Getenv("PORT")
	server := NewAPIServer("0.0.0.0:"+port)
	server.Serve()
}

func NewAPIServer(listenAdd string) *APIServer {
	return &APIServer{
		listenAdd: listenAdd,
	}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET"})
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Authorization"})

	router.HandleFunc("/history", handelHistory)

	log.Println("Listening on " + s.listenAdd)
	log.Fatal(http.ListenAndServe(s.listenAdd, handlers.CORS(credentials, origins, methods, headers)(router)))
}

func handelHistory(w http.ResponseWriter, r *http.Request) {
	// validate github user
	repo := r.URL.Query().Get("repo")
	token := r.Header.Get("Authorization")

	fmt.Println("client token: ", token)

	if repo == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res := map[string]string{
			"message": "repository name required",
			"code":    "400",
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	if token == " " {
		token = fmt.Sprintf("Bearer %s", os.Getenv("GITHUB_TOKEN"))
	}

	fmt.Println(token)

	owner := strings.Split(repo, "/")[0]
	logoUrl, err := api.GetRepoLogoUrl(owner, token)
	if err != nil {
		returnErr(w, err.Code, err.Message)
		return
	}

	// validate repository
	totalStars, err := api.GetRepoTotalStarCount(repo, token)
	if err != nil {
		returnErr(w, err.Code, err.Message)
		return
	}

	totalPageCount, err := api.GetRepoPageCount(repo, token)
	if err != nil {
		returnErr(w, err.Code, err.Message)
		return
	}

	// collect stars record
	record, err := api.GetRepoStargazers(repo, token, 15, totalPageCount, totalStars)

	if err != nil {
		returnErr(w, err.Code, err.Message)
		return
	}

	response := map[string]interface{}{
		"total_stars": totalStars,
		"logo_url":    logoUrl,
		"data": 	  record,
		"name":        repo,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	
}

func returnErr(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	res := map[string]string{
		"message": message,
		"code":    strconv.Itoa(code),
	}
	json.NewEncoder(w).Encode(res)
}
