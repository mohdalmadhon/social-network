package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social/sql/database"
)

func (app *App) GetUserData(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := database.GetUserData(app.DB, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	profileData, err := database.GetProfileData(app.DB, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get profile data",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": true,
		"data": map[string]any{
			"user":    data,
			"profile": profileData,
		},
	})
}
