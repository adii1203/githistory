package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"src/api"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAdd string
}

func NewAPIServer(listenAdd string) *APIServer {
	return &APIServer{
		listenAdd: listenAdd,
	}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/history", handelHistory)

	log.Println("Listening on " + s.listenAdd)
	log.Fatal(http.ListenAndServe(s.listenAdd, router))
}

func main() {
	port := os.Getenv("PORT")
	server := NewAPIServer(":"+port)
	server.Serve()
}

func handelHistory(w http.ResponseWriter, r *http.Request) {
	// validate github user
	repo := r.URL.Query().Get("repo")
	token := r.Header.Get("Authorization")

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
	owner := strings.Split(repo, "/")[0]
	_, err := api.GetRepoLogoUrl(owner, token)
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
	record := api.GetRepoStargazers(repo, token, 15, totalPageCount, totalStars)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
	return
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
