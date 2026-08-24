package api

import (
	"encoding/json"
	"net/http"
	"social/sql/database"
)

func (app *App) GetUserData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "bad request method",
		})
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get cookie",
		})
		return
	}

	payload, err := VerifyToken(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "unauthorized",
		})
		return
	}

	data, err := database.GetUserDataByIdentifier(app.DB, payload.Identifier)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	profileData, err := database.GetProfileData(app.DB, payload.Identifier)
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
