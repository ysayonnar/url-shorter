package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	customErrors "url-shorter/pkg/errors"
)

type UrlShortedResponse struct{
	InitialUrl string `json:"initialUrl"`
	ShortedUrl string `json:"shortedUrl"`
}

func IsUrlExists(url string) (bool, customErrors.DefaultError){
	client := &http.Client{
		Timeout: 2 * time.Second,
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

func SendJsonResponse(w http.ResponseWriter, r *http.Request, response  UrlShortedResponse){
	jsonResponse, parsingErr := json.Marshal(response)
	if parsingErr != nil {
		err := customErrors.DefaultError{
			Message: "Error while parsing response object",
			StatusCode: http.StatusInternalServerError,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}