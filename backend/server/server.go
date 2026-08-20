package server

import (
	"database/sql"
	"net/http"
	"social/backend/server/api"
)


func StartServer(db *sql.DB) *http.ServeMux {
	app := api.App{
		DB: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/register/checkEmail", app.CheckEmailExists)
	mux.HandleFunc("/api/register/checkUsername", app.CheckUsernameExists)
	mux.HandleFunc("/api/register/submit", app.RegisterUser)
	mux.HandleFunc("/api/login", app.LoggingUser)
	return mux
}
