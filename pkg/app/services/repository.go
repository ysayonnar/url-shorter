package services

import (
	"database/sql"
	"log"
)

func InsertUrl(db *sql.DB, oldUrl, newUrl string) error{
	query := `
        INSERT INTO urls (oldurl, newurl)
        VALUES ($1, $2)
    `
	_, err := db.Exec(query, oldUrl, newUrl)
	if err != nil{
		log.Println("Error while inserting data", err)
		return err
	}
	log.Println("Successfully inserted!")
	return nil
}