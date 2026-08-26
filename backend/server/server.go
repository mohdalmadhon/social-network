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
	mux.HandleFunc("DELETE /api/session", app.LogOutUser)
	//users
	mux.HandleFunc("POST /api/user/", app.RegisterUser)
	mux.HandleFunc("PUT /api/user", middleware.AuthMiddleware(app.UpdateUserInfo))
	mux.HandleFunc("PUT /api/user/avatar", middleware.AuthMiddleware(app.UpdateUserAvatar))
	mux.HandleFunc("DELETE /api/user/avatar", middleware.AuthMiddleware(app.DeleteuserAvatar))
	mux.HandleFunc("GET /api/user/check-email", app.CheckEmailExists)
	mux.HandleFunc("GET /api/user/check-username", app.CheckUsernameExists)

	//profile
	mux.HandleFunc("GET /api/me", middleware.AuthMiddleware(app.GetUserData))

	mux.HandleFunc("/api/posts", app.Posts)

	uploadDir, err := filepath.Abs("../uploads")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle(
		"/uploads/",
		http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))),
	)

	mux.Handle(
		"/images/",
		http.StripPrefix("/images/", http.FileServer(http.Dir("../images"))),
	)
	return mux
}
