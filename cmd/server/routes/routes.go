package routes

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"social/internal/app/api"
)

func StartServer(db *sql.DB) *http.ServeMux {
	uploadsDir, err := filepath.Abs("uploads")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	app := api.App{DB: db}

	mux.HandleFunc("POST /api/user", app.RegisterUser)
	mux.Handle("/uploads/",http.StripPrefix("/uploads/",http.FileServer(http.Dir(uploadsDir)),),)
	return mux
}
