package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"social/database/users"
	"social/internal/app/tokens"
)
/*
A middle ware so simply takes a handler and returns a handler

In this handler we check if the user is authoniticated by checking the cookies which use JWT format. then Decode and check the signature,
and Deconde and check the user ID. then write the user ID in the request Context for the callback to use it.

Paramters:
	handler http.HandleFunc

Returns:
	http.HandleFunc
*/
func (app *App) AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "not authenticated",
			})
			return
		}

		payload, err := tokens.VerifyToken(cookie.Value)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "invalid or expired session",
			})
			return
		}

		err = users.UserExists(app.DB, payload.UserID)
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "invalid user",
			})
			return
		}

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "could not verify user",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "userID", payload.UserID)

		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}