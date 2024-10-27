package handlers

import (
	"github.com/gorilla/mux"
)

func SetHandlers() *mux.Router {
	router := mux.NewRouter()
	
	router.HandleFunc("/api/v1/short", GenerateShortedUrl)
	return router
}