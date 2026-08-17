package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/social_network.db")
	if err != nil {
		log.Printf("Error opening to a database %s", err.Error())
		return
	}
	defer db.Close()

	err = db.Ping() 
	if err != nil {
		log.Printf("Error connecting to the database: %s", err.Error())
		return
	}
	
	log.Println("Successfully connected to the database")
}