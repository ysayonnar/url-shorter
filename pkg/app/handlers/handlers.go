package handlers

import (
	"fmt"
	"net/http"
	customErrors "url-shorter/pkg/errors"
)

func GenerateShortedUrl(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		err := customErrors.DefaultError{Message: "Method must be POST", StatusCode: http.StatusMethodNotAllowed}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}
	fmt.Fprint(w, "Hello, world!")
}