package customErrors

import (
	"encoding/json"
	"log"
	"net/http"
)

type DefaultError struct{
	Message string `json:"msg"`
	StatusCode int `json:"-"`
}

func ThrowDefaultError(w http.ResponseWriter, r *http.Request, e DefaultError){
	jsonErr, err := json.Marshal(e)
	if err != nil {
		log.Fatal("Error while parsing <DefaultError> to json.")
	}
	w.WriteHeader(e.StatusCode)
	w.Write(jsonErr)
}