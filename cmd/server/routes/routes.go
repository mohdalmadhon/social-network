package routes

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"social/internal/app/api"
	middleware "social/internal/app/middleWare"
)

func StartServer(db *sql.DB) *http.ServeMux {
	uploadsDir, err := filepath.Abs("uploads")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	app := api.App{DB: db}

	//user
	mux.HandleFunc("GET /api/user", middleware.AuthMiddleware(app.GetUserData))
	mux.HandleFunc("POST /api/user", app.RegisterUser)
	mux.HandleFunc("PATCH /api/user", middleware.AuthMiddleware(app.UpdateUserInfo))

	//session
	mux.HandleFunc("POST /api/session", app.LoggingUser)
	mux.HandleFunc("GET /api/session", app.AuthorizeSession)
	mux.HandleFunc("DELETE /api/session", app.DeleteSession)
	
	//profiles
	mux.HandleFunc("PATCH /api/profile/avatar", middleware.AuthMiddleware(app.UpdateUserAvatar))
	mux.HandleFunc("GET /api/profile/about", middleware.AuthMiddleware(app.GetUserAbout))
	mux.HandleFunc("PATCH /api/profile/about", middleware.AuthMiddleware(app.UpdateUserAbout))
	mux.HandleFunc("GET /api/profile", middleware.AuthMiddleware(app.GetUserProfile))

	// follow handler
	mux.HandleFunc("POST /api/profile/follow", middleware.AuthMiddleware(app.RequestFollow))
	mux.HandleFunc("DELETE /api/profile/follow", middleware.AuthMiddleware(app.CancelRequest))
	mux.HandleFunc("GET /api/profile/follow", middleware.AuthMiddleware(app.GetFollowers))
	mux.HandleFunc("GET /api/profile/following", middleware.AuthMiddleware(app.GetFollowing))
	
	//folder handlers
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))
	return mux
}
