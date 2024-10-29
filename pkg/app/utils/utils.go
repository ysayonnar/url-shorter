package utils

import (
	"log"
	"net/http"
	"time"
	customErrors "url-shorter/pkg/errors"
)

func IsUrlExists(url string) (bool, customErrors.DefaultError){
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Head(url)
	if err != nil{
		log.Fatal(err)
		customErr := customErrors.DefaultError{
			Message: "Error while checking is url exist",
			StatusCode: http.StatusInternalServerError,
		}
		return false, customErr
	}
	defer response.Body.Close()

	return response.StatusCode >= 200 && response.StatusCode < 400, customErrors.DefaultError{}
}