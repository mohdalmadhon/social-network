package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"social/internal/app/tokens"
)

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
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

		ctx := context.WithValue(r.Context(), "userID", payload.UserID)

		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}
