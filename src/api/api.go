package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"sort"
	"src/utils"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type StarRecord struct {
	Date  string `json:"date"`
	Stars int    `json:"stars"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}

func GetRepoLogoUrl(owner string, token string) (string, *ErrorResponse) {
	client := &http.Client{}
	url := "https://api.github.com/users/" + owner
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Accept", "application/vnd.github.v3+json")
	request.Header.Add("Authorization", token)
	response, err := client.Do(request)
	statusCode := response.StatusCode

	if statusCode == 403 {
		return "", &ErrorResponse{
			Code:    403,
			Message: "rate limit exceeded",
		}
	} else if statusCode == 500 {
		return "", &ErrorResponse{
			Code:    500,
			Message: "server error",
		}
	} else if statusCode == 404 {
		return "", &ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("user %s not found", owner),
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)

	jsonData := map[string]interface{}{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData["avatar_url"].(string), nil
}

func GetRepoTotalStarCount(repo, token string) (int, *ErrorResponse) {
	client := &http.Client{}
	url := "https://api.github.com/repos/" + repo
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Accept", "application/vnd.github.v3.star+json")
	request.Header.Add("Authorization", token)

	response, err := client.Do(request)

	statusCode := response.StatusCode

	if statusCode == 403 {
		return 0, &ErrorResponse{
			Code:    403,
			Message: "rate limit exceeded",
		}
	} else if statusCode == 500 {
		return 0, &ErrorResponse{
			Code:    500,
			Message: "server error",
		}
	} else if statusCode == 404 {
		return 0, &ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("repository %s not found", repo),
		}
	} else if err != nil {
		return 0, &ErrorResponse{
			Code:    500,
			Message: err.Error(),
		}
	}

	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)
	jsonData := map[string]interface{}{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	if int(jsonData["stargazers_count"].(float64)) == 0 {
		return 0, &ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("repo %s dose not have any stars", repo),
		}
	}
	return int(jsonData["stargazers_count"].(float64)), nil
}

func GetRepoPageCount(repo, token string) (int, *ErrorResponse) {
	client := &http.Client{}
	url := "https://api.github.com/repos/" + repo + "/stargazers"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Add("Authorization", token)

	response, err := client.Do(request)
	statusCode := response.StatusCode

	if statusCode == 403 {
		return 0, &ErrorResponse{
			Code:    403,
			Message: "rate limit exceeded",
		}
	} else if statusCode == 500 {
		return 0, &ErrorResponse{
			Code:    500,
			Message: "server error",
		}
	} else if statusCode == 404 {
		return 0, &ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("repository %s not found", repo),
		}
	} else if err != nil {
		return 0, &ErrorResponse{
			Code:    500,
			Message: err.Error(),
		}
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	link := response.Header["Link"]

	if link == nil && string(body) != "[]" {
		return 0, nil
	}

	l := link[0]
	totalPages := 1
	regex := regexp.MustCompile("<[^>]*?page=(\\d+)>; rel=\"last\"")
	regexResult := regex.FindStringSubmatch(l)
	if len(regexResult) != 0 {
		pageCount := regexResult[1]
		if num, err := strconv.Atoi(pageCount); err == nil {
			totalPages = num
		}
	}

	return totalPages, nil
}

// GetStargazersCountPerPage todo: refactor
func GetStargazersCountPerPage(repo, token string, page int) (map[string]interface{}, error) {
	client := &http.Client{}

	url := "https://api.github.com/repos/" + repo + "/stargazers?per_page=30" + "&page=" + strconv.Itoa(page)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Accept", "application/vnd.github.v3.star+json")
	request.Header.Add("Authorization", token)

	response, _ := client.Do(request)

	if response.StatusCode != 200 {
		e := fmt.Errorf("something went wrong, status: %d", response.StatusCode)
		return nil, e
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var parsedData []map[string]interface{}

	if string(body) == "[]" {
		return make(map[string]interface{}), nil
	}

	data := map[string]interface{}{
		"data": []map[string]string{},
	}

	err = json.Unmarshal(body, &parsedData)
	if err != nil {
		log.Fatal(err)

	}
	for i := 0; i < len(parsedData); i++ {
		tt := parsedData[i]["starred_at"].(string)
		data["data"] = append(data["data"].([]map[string]string), map[string]string{
			"starred_at": tt,
		})
	}

	return data, nil
}

// GetRepoStargazers todo: refactor
func GetRepoStargazers(repo, token string, maxRequestAmount int, totalPage, totalStars int) []StarRecord {

	var totalRequestPages []int
	if totalPage < maxRequestAmount {
		totalRequestPages = append(totalRequestPages, utils.GetRange(1, totalPage)...)
	} else {
		for _, value := range utils.GetRange(1, maxRequestAmount) {
			totalRequestPages = append(totalRequestPages, int(math.Round(float64(value*totalPage/maxRequestAmount))-1))
		}
		if totalRequestPages[0] != 1 {
			totalRequestPages[0] = 1
		}
	}

	resArray := make([][]map[string]string, len(totalRequestPages))

	ch := make(chan []map[string]string, len(totalRequestPages))
	for _, value := range totalRequestPages {
		go func(v int) {
			res, err := GetStargazersCountPerPage(repo, token, v)
			if err != nil {
				log.Fatal(err)
			}
			ch <- res["data"].([]map[string]string)
		}(value)
		//resArray[idx] = v["data"].([]map[string]string)
	}

	for i := 0; i < len(totalRequestPages); i++ {
		resArray[i] = <-ch
	}

	starRecordMap := make(map[string]int)

	if len(totalRequestPages) < maxRequestAmount {
		var starRecordData []map[string]string
		for _, res := range resArray {
			starRecordData = append(starRecordData, res...)
		}

		for i := 0; i < len(starRecordData); {
			starRecordMap[utils.GetTimeStamp(starRecordData[i]["starred_at"])] = i + 1
			step := int(math.Floor(float64(len(starRecordData) / maxRequestAmount)))
			if step == 0 {
				step = 1
			}
			i += step
		}
	} else {
		for idx, data := range resArray {
			if len(data) > 0 {
				starRecord := data[0]
				starRecordMap[utils.GetTimeStamp(starRecord["starred_at"])] = 30 * (totalRequestPages[idx] - 1)
			}
		}
	}

	currentDate := time.Now()
	parsedDate, _ := time.Parse("2006-01-02 15:04:05.999999999", currentDate.Format("2006-01-02 15:04:05.999999999"))
	starRecordMap[parsedDate.Format("Jan 2006")] = totalStars

	var starRecords []StarRecord

	for key, value := range starRecordMap {
		starRecords = append(starRecords, StarRecord{
			Stars: value,
			Date:  key,
		})
	}

	sort.Slice(starRecords, func(i, j int) bool {
		return starRecords[i].Stars < starRecords[j].Stars
	})
	return starRecords
}
