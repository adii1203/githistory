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
	"sync"
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
		return "", &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to create request: %v", err),
		}
	}
	request.Header.Add("Accept", "application/vnd.github.v3+json")
	request.Header.Add("Authorization", token)
	response, err := client.Do(request)
	if err != nil {
		return "", &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to make request: %v", err),
		}
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 403:
		return "", &ErrorResponse{
			Code:    403,
			Message: "rate limit exceeded",
		}

	case 500:
		return "", &ErrorResponse{
			Code:    500,
			Message: "server error",
		}
	case 404:
		return "", &ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("user %s not found", owner),
		}
	case 200:

	case 401:
		return "", &ErrorResponse{
			Code:    401,
			Message: "Bad credentials",
		}
	default:
		return "", &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("unexpected error: %v", response.Status),
		}

	}



	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to read response body: %v", err),
		}
	}

	jsonData := map[string]interface{}{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return "", &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to parse response body: %v", err),
		}
	}
	return jsonData["avatar_url"].(string), nil
}

func GetRepoTotalStarCount(repo, token string) (int, *ErrorResponse) {
	client := &http.Client{}
	url := "https://api.github.com/repos/" + repo
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return 0, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to create request: %v", err),
		}
	}
	request.Header.Set("Accept", "application/vnd.github.v3.star+json")
	request.Header.Add("Authorization", token)

	response, err := client.Do(request)
	if err != nil {
		return 0, &ErrorResponse{
			Code:   500,
			Message: fmt.Sprintf("faild to make request: %v", err),
		}
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 403:
		return 0, &ErrorResponse{
			Code:    403,
			Message: "rate limit exceeded",
		}

	case 500:
		return 0, &ErrorResponse{
			Code:    500,
			Message: "server error",
		}
	case 404:
		return 0, &ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("repository %s not found", repo),
		}
	case 200:

	case 401:
		return 0, &ErrorResponse{
			Code:    401,
			Message: "Bad credentials",
		}
	default:
		return 0, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("unexpected error: %v", response.Status),
		}

	}
	
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to read response body: %v", err),
		}
	}

	jsonData := map[string]interface{}{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return 0, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to parse response body: %v", err),
		}
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to read response body: %v", err),
		}
	}

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

	url := fmt.Sprintf("https://api.github.com/repos/%s/stargazers?per_page=30&page=%d", repo, page)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to create request: %v", err),
		}
	}

	request.Header.Set("Accept", "application/vnd.github.v3.star+json")
	request.Header.Add("Authorization", token)

	response, err := client.Do(request)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to make request: %v", err),
		}
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 403:
		return nil, &ErrorResponse{
			Code:    403,
			Message: "rate limit exceeded",
		}
	case 500:
		return nil, &ErrorResponse{
			Code:    500,
			Message: "server error",
		}
	case 404:
		return nil, &ErrorResponse{
			Code:    404,
			Message: fmt.Sprintf("repository %s not found", repo),
		}
	case 200:

	default:
		return nil, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("unexpected error: %v", response.Status),
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to read response body: %v", err),
		}
	}

	if string(body) == "[]"{
		return map[string]interface{}{"data": []map[string]string{}}, nil
	}

	var parsedData []map[string]interface{}
	err = json.Unmarshal(body, &parsedData)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("faild to parse response body: %v", err),
		}
	}

	data := map[string]interface{}{
		"data": []map[string]string{},
	}
	for _, item := range parsedData{
		if starrdAt, ok := item["starred_at"].(string); ok {
			data["data"] = append(data["data"].([]map[string]string), map[string]string{
				"starred_at": starrdAt,
			})
		}
	}

	return data, nil
}

func GetRepoStargazers(repo, token string, maxRequestAmount int, totalPage, totalStars int) ([]StarRecord, *ErrorResponse) {

	totalRequestPages := calculateRequestPages(maxRequestAmount, totalPage)

	resArray := make([][]map[string]string, len(totalRequestPages))
	errorList := []error{}

	var mu sync.Mutex
	var wg sync.WaitGroup
	for idx, page := range totalRequestPages {
		wg.Add(1)
		go func(i, p int) {
			defer wg.Done()
			res, err := GetStargazersCountPerPage(repo, token, p)
			if err != nil {
				mu.Lock()
				errorList = append(errorList, err)
				mu.Unlock()
				return
			}
			
			if data, ok := res["data"].([]map[string]string); ok {
				mu.Lock()
				resArray[i] = data
				mu.Unlock()
			}
		}(idx, page)
	}

	wg.Wait()

	if len(errorList) > 0 {
		return nil, &ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("error occured: %v", errorList),
		}
	}

	starRecords := aggregateStarData(resArray, totalRequestPages, maxRequestAmount, totalStars)

	sort.Slice(starRecords, func(i, j int) bool {
		return starRecords[i].Stars < starRecords[j].Stars
	})
	return starRecords, nil
}

func aggregateStarData(resArray [][]map[string]string, totalRequestPages []int, maxRequestAmount, totalStars int) []StarRecord {
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

	return starRecords
}

func calculateRequestPages(maxRequestAmount, totalPage int) []int{
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
	return totalRequestPages
}