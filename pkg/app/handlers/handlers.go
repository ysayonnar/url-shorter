package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"url-shorter/pkg/app/services"
	"url-shorter/pkg/app/utils"
	"url-shorter/pkg/config"
	customErrors "url-shorter/pkg/errors"

	"github.com/gorilla/mux"
)

type Body struct{
	Url string `json:"url"`
	Length int `json:"length"`
}

type ClicksJson struct{
	OldUrl string `json:"oldurl"`
	NewUrl string `json:"newurl"`
	Clicks int `json:"clicks"`
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func GenerateShortedUrl(w http.ResponseWriter, r *http.Request, db *sql.DB){
	if r.Method != "POST"{
		err := customErrors.DefaultError{Message: "Method must be POST", StatusCode: http.StatusMethodNotAllowed}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}

	log.Println(Green + "POST URL:" + Reset + r.URL.String())

	body := &Body{}
	parsingErr := json.NewDecoder(r.Body).Decode(body)
	if parsingErr != nil || len(body.Url) == 0{
		err := customErrors.DefaultError{
			Message: "Incorrect body", 
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}
	if body.Length < 8 || body.Length > 20{
		err := customErrors.DefaultError{
			Message: "Length must be from 8 to 20",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w,r,err)	
		return
	}

	if utils.CheckIsUrlShorted(body.Url){
		err := customErrors.DefaultError{
			Message: "Url is already shorted and cant be shorted again",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w,r,err)
		return
	}
	
	IsUrlExists, err := utils.IsUrlExists(body.Url)
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
	// end of checking url existing


	//check is created already
	urlRecord, checkExistError := services.GetUrlByOldUrl(db, body.Url)
	if checkExistError != nil{
		customErrors.ThrowDefaultError(w, r, *checkExistError)
		return
	}
	if urlRecord != nil {
		response := utils.UrlShortedResponse{
			InitialUrl: body.Url,
			ShortedUrl: urlRecord.Newurl,
		}
		utils.SendJsonResponse(w, r, response)
		return
	}

	//genearting token url
	shortedUrl, err := services.UrlShorter(body.Url, body.Length)
	if err.Message != ""{
		customErrors.ThrowDefaultError(w,r, err)
		return
	}

	queryErr := services.InsertUrl(db, body.Url, shortedUrl)
	if queryErr != nil{
		err := customErrors.DefaultError{
			Message: err.Message,
			StatusCode: http.StatusInternalServerError,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}
	
	response := utils.UrlShortedResponse{
		InitialUrl: body.Url,
		ShortedUrl: shortedUrl,
	}
	utils.SendJsonResponse(w, r, response)
}

func Redirect(w http.ResponseWriter, r *http.Request, db *sql.DB){
	//ТУТ СДЕЛАТЬ ДОБАВЛЕНИЕ КЛИКОВ ПРИ ПЕРЕХОДЕ
	if r.Method != "GET"{
		err := customErrors.DefaultError{
			Message: "Method must be GET",
			StatusCode: http.StatusMethodNotAllowed,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}

	log.Println(Green + "GET" + Reset + r.URL.String())

	params := mux.Vars(r)
	token, isExists := params["token"]
	if !isExists{
		err := customErrors.DefaultError{
			Message: "Token is invalid",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}
	_, err := strconv.Atoi(token)
	if err == nil{
		err := customErrors.DefaultError{
			Message: "Token cant be a number",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}

	newUrl := fmt.Sprintf("https://%v/%s", config.Domen, token) //generating url

	record, queryError := services.GetUrlByNewUrl(db, newUrl)
	if queryError != nil {
		customErrors.ThrowDefaultError(w, r, *queryError)
		return
	}
	if record == nil{
		//getting current working directory
		cwd, _ := os.Getwd()
		cwd = filepath.Clean(cwd)
		path := filepath.Join(cwd, "../pkg/htmlPages/index.html")

		htmlFile, err := os.ReadFile(path)
		if err != nil{
			log.Println("Error while getting html")
			log.Fatal(err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlFile)
		return
	}
	newClicks := record.Clicks + 1
	clicksErr := services.IncreaseClicks(db, newClicks, record.Id)
	if clicksErr != nil{
		customErrors.ThrowDefaultError(w, r, *clicksErr)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Redirect(w,r, record.Oldurl, http.StatusFound)
}

func GetClicks(w http.ResponseWriter, r *http.Request, db *sql.DB){
	if r.Method != "GET"{
		err := customErrors.DefaultError{
			Message: "Method must be GET",
			StatusCode: http.StatusMethodNotAllowed,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}

	log.Println(Green + "GET" + Reset + r.URL.String())

	params := mux.Vars(r)
	token, isExists := params["token"]
	if !isExists{
		err := customErrors.DefaultError{
			Message: "Token is invalid",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}
	_, err := strconv.Atoi(token)
	if err == nil{
		err := customErrors.DefaultError{
			Message: "Token cant be a number",
			StatusCode: http.StatusBadRequest,
		}
		customErrors.ThrowDefaultError(w, r, err)
		return
	}

	newUrl := fmt.Sprintf("https://%v/%s", config.Domen, token) //generating url

	record, queryError := services.GetUrlByNewUrl(db, newUrl)
	if queryError != nil {
		customErrors.ThrowDefaultError(w, r, *queryError)
		return
	}
	if record == nil{
		//getting current working directory
		cwd, _ := os.Getwd()
		cwd = filepath.Clean(cwd)
		path := filepath.Join(cwd, "../pkg/htmlPages/index.html")

		htmlFile, err := os.ReadFile(path)
		if err != nil{
			log.Println("Error while getting html")
			log.Fatal(err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlFile)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	clicksJson := ClicksJson{
		Clicks: record.Clicks,
		OldUrl: record.Oldurl,
		NewUrl: record.Newurl,
	}
	jsonData, err := json.Marshal(clicksJson)
	if err != nil{
		customErrors.ThrowDefaultError(w, r, customErrors.DefaultError{
			Message: "Error whileparsing json",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	w.Write(jsonData)
}