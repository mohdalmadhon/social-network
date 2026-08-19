package server

import (
	"database/sql"
	"net/http"
	"social/backend/server/api"
)


func startServer(db *sql.DB) *http.ServeMux {
	app := api.App{
		DB: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/checkEmail", app.CheckEmailExists)
	
	return mux
}
