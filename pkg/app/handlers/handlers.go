package handlers

import (
	"encoding/json"
	"net/http"
	"url-shorter/pkg/app/services"
	"url-shorter/pkg/app/utils"
	customErrors "url-shorter/pkg/errors"
)

type Url struct{
	Url string `json:"url"`
	Length int `json:"length"`
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
		err := customErrors.DefaultError{
			Message: "Incorrect body", 
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}
	if url.Length < 8 || url.Length > 20{
		err := customErrors.DefaultError{
			Message: "Length must be from 8 to 20",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w,r,err)	
		return
	}

	//1.проверка на то что ссылка не моя
	//2. проверка создана ли уже ссылка, если да, то просто доставать из базы и кидать обратно
	//3. проверка на существование ссылки
	
	IsUrlExists, err := utils.IsUrlExists(url.Url)
	if err.Message != ""{
		customErrors.ThrowDefaultError(w,r,err)
		return
	}
	
	if !IsUrlExists{
		err := customErrors.DefaultError{
			Message: "Url does not exits",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w,r, err)
		return
	}
	
	shortedUrl, err := services.UrlShorter(url.Url, url.Length)
	if err.Message != ""{
		customErrors.ThrowDefaultError(w,r, err)
		return
	}
	// end of checking url existing
	
	response := UrlShortedResponse{
		InitialUrl: url.Url,
		ShortedUrl: shortedUrl,
	}
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