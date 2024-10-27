package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	customErrors "url-shorter/pkg/errors"
)

type Url struct{
	Url string `json:"url"`
}

func GenerateShortedUrl(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		err := customErrors.DefaultError{Message: "Method must be POST", StatusCode: http.StatusMethodNotAllowed}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}

	url := &Url{}
	parsingErr := json.NewDecoder(r.Body).Decode(url)
	if parsingErr != nil{
		err := customErrors.DefaultError{Message: "Incorrect body", StatusCode: http.StatusBadRequest}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}
	
	fmt.Fprintf(w, url.Url)
}