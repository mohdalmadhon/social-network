package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"social/backend/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "../db/social_network.db")
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

	s := http.Server{
		Handler: server.StartServer(db),
		Addr:    ":4033",
	}

	fmt.Println("Server running on: http://localhost" + s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal("Error listening to server", err)
	}
}
