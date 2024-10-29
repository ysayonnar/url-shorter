package customErrors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type DefaultError struct{
	Message string `json:"msg"`
	StatusCode int `json:"-"`
}
func (e *DefaultError) Error() string{
	return fmt.Sprintf("Code: %v, Message: %s", e.StatusCode, e.Message)
}

func ThrowDefaultError(w http.ResponseWriter, r *http.Request, e DefaultError){
	jsonErr, err := json.Marshal(e)
	if err != nil {
		log.Fatal("Error while parsing <DefaultError> to json.")
	}
	w.WriteHeader(e.StatusCode)
	w.Write(jsonErr)
}
