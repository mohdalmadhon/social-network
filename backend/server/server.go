package server

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"social/backend/server/api"
	middleware "social/backend/server/api/middleWare"
)

func StartServer(db *sql.DB) *http.ServeMux {
	app := api.App{
		DB: db,
	}

	mux := http.NewServeMux()
	// sessions
	mux.HandleFunc("POST /api/session", app.LoggingUser)

	//users
	mux.HandleFunc("POST /api/user/", app.RegisterUser)
	mux.HandleFunc("PUT /api/user", middleware.AuthMiddleware(app.UpdateUserInfo))
	
	//profile
	mux.HandleFunc("GET /api/me", middleware.AuthMiddleware(app.GetUserData))

	mux.HandleFunc("/api/register/checkEmail", app.CheckEmailExists)
	mux.HandleFunc("/api/register/checkUsername", app.CheckUsernameExists)
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
