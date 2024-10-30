package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"url-shorter/pkg/app/database"
	"url-shorter/pkg/app/handlers"
)

func main() {
	// log.SetFlags(log.Ltime) //leave only the time in the logs
	args := os.Args
	//checking arguments and port
	if len(args) != 2{
		log.Fatal("Incorrect arguments, usage: main.exe <port>")
		return
	}
	port, err := strconv.Atoi(args[1])
	if err != nil{
		log.Fatal("<port> must be an integer.")
		return
	}
	if port < 1000 || port > 9999{
		log.Fatal("<port> must be from 1000 to 9999.")
	}
	
	db, err := database.ConnectDb()
	if err != nil{
		return
	}
	err = database.CreateTables(db)
	if err != nil{
		return
	}

	log.Printf("Server started on port %d", port);
	http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.SetHandlers(db))
}