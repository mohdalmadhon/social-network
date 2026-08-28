package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/database/users"
	"social/internal/app/tokens"
	"social/internal/helpers"
	"social/internal/models"
	"time"
)

func (app *App) LoggingUser(w http.ResponseWriter, r *http.Request) {
	var logger models.UserLogger

	if err := json.NewDecoder(r.Body).Decode(&logger); err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "could not get data",
		})
		return
	}

	hashedPassword, err := database.GetHashedPassowrd(app.DB, logger.Identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
				"status":  false,
				"message": "user identifier or password is incorrect",
			})
			return
		} else {
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not get user data",
			})
			return
		}
	}

	if match := helpers.AuthonticateUser(logger.Pass, hashedPassword); !match {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "user identifier or password is incorrect",
		})
		return
	}

	userID := database.GetUserID(app.DB, logger.Identifier)
	if userID == -1 {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}
	
	token, err := tokens.GenerateToken(userID)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not create authentication token",
		})
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * 30 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "logged in",
	})
}
