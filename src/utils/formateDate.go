package utils

import (
	"log"
	"time"
)

func GetTimeStamp(t string) string {
	value, err := time.Parse(time.RFC3339, t)
	if err != nil {
		log.Fatal(err)
	}
	formatedTime := value.Format("Jan 2006")

	return formatedTime
}
