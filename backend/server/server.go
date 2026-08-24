package server

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
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
	mux.HandleFunc("/api/session", app.SessionAuthorizer)
	mux.HandleFunc("/api/users/getData", app.GetUserData)
	mux.HandleFunc("/api/posts", app.Posts)

	uploadDir, err := filepath.Abs("../uploads")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle(
		"/uploads/",
		http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))),
	)
	return mux
}
