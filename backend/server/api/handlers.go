package api

import (
	"database/sql"
	"net/http"
)

type App struct {
	DB *sql.DB
}

func (app App) CheckEmailExists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// error page
		return
	}
}
