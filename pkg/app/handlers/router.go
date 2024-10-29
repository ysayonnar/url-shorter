package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func SetHandlers(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	
	router.HandleFunc("/short/{token}", func(w http.ResponseWriter, r *http.Request) {
		Redirect(w,r,db)
	})
	router.HandleFunc("/api/v1/short", func(w http.ResponseWriter, r *http.Request) {
		GenerateShortedUrl(w, r, db)
	})
	return router
}