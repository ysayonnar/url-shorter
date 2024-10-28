package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"url-shorter/pkg/app/services"
	customErrors "url-shorter/pkg/errors"
)

type Url struct{
	Url string `json:"url"`
}

type UrlShortedResponse struct{
	InitialUrl string `json:"initialUrl"`
	ShortedUrl string `json:"shortedUrl"`
}

func GenerateShortedUrl(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		err := customErrors.DefaultError{Message: "Method must be POST", StatusCode: http.StatusMethodNotAllowed}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}

	url := &Url{}
	parsingErr := json.NewDecoder(r.Body).Decode(url)
	if parsingErr != nil || len(url.Url) == 0{
		err := customErrors.DefaultError{Message: "Incorrect body", StatusCode: http.StatusBadRequest}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}
	
	shortedUrl := services.UrlShorter(url.Url)//тут сделать возможность обработки ошибки типо DefaultError из customErrors
	//сделать проверку на то чтобы ссылки не было в бд
	//если ссылка существует то просто ее же и возвращать
	// сделать проверку на существование ссылки
	
	response := UrlShortedResponse{
		InitialUrl: url.Url,
		ShortedUrl: shortedUrl,
	}
	jsonResponse, parsingErr := json.Marshal(response)
	if parsingErr != nil {
		log.Fatal("Error while parsing", parsingErr)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}