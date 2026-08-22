package api

import (
	"encoding/json"
	"net/http"
)

func (app App) SessionAuthorizer(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get cookie",
		})
		return
	}

	payload, err := VerifyToken(token.Value)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get cookie",
		})
		return
	}
	if payload == nil {
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get cookie",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "valid token",
	})
	return
}
