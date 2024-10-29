package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDb()(*sql.DB, error){
	connectinStr := "user=postgres password=root dbname=url_shorter sslmode=disable"
	db, err := sql.Open("postgres", connectinStr)
	if err != nil {
		log.Println("Error while connection to database.")
		log.Fatal(err)
		return nil, err
	}
	// defer db.Close()
	
	err = db.Ping()
	if err != nil{
		log.Println("Erorr while pinging database")
		log.Fatal(err)
		return nil, err
	}

	log.Println("Connected to database")
	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
				id SERIAL PRIMARY KEY,
				oldUrl TEXT NOT NULL,
				newUrl TEXT NOT NULL,
				clicks INTEGER DEFAULT 0
			);`
	_, err := db.Exec(query)
	if err != nil{
		log.Fatalf("Error while creating tables: %w", err)
		return err
	}
	log.Println("Tables created")
	return nil
}