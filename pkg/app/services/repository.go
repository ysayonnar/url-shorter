package services

import (
	"database/sql"
	"log"
	"net/http"
	customErrors "url-shorter/pkg/errors"
)

type UrlRecord struct{
	Id int
	Oldurl string
	Newurl string
	Clicks int
}

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
	return nil
}

func GetUrlByOldUrl(db *sql.DB, oldUrl string) (*UrlRecord, *customErrors.DefaultError){
	query := "SELECT * FROM urls WHERE oldurl = $1"
	var urlRecord UrlRecord
	err := db.QueryRow(query, oldUrl).Scan(&urlRecord.Id, &urlRecord.Oldurl, &urlRecord.Newurl, &urlRecord.Clicks)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil{
		log.Println("Error while connecting to database")
		return nil, &customErrors.DefaultError{Message: "Error while connecting to database", StatusCode: http.StatusInternalServerError}
	}
	return &urlRecord, nil
}

func GetUrlByNewUrl(db *sql.DB, newUrl string) (*UrlRecord, *customErrors.DefaultError){
	query := "SELECT * FROM urls WHERE newurl = $1"
	var urlRecord UrlRecord
	err := db.QueryRow(query, newUrl).Scan(&urlRecord.Id, &urlRecord.Oldurl, &urlRecord.Newurl, &urlRecord.Clicks)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil{
		log.Println("Error while connecting to database")
		return nil, &customErrors.DefaultError{Message: "Error while connecting to database", StatusCode: http.StatusInternalServerError}
	}
	return &urlRecord, nil
}

func IncreaseClicks(db *sql.DB, clicks int, id int) *customErrors.DefaultError{
	query := `UPDATE urls SET clicks = $1 WHERE id = $2`
	_, err := db.Exec(query, clicks, id)
	if err != nil{
		log.Println("Error while updating the data!")
		return  &customErrors.DefaultError{
			Message: "Error while updating the data",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}