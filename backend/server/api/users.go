package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"social/backend/server/api/helpers"
	"social/sql/database"

	"github.com/google/uuid"
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

func (app *App) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		if err == http.ErrMissingFile {
			helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "no image uploaded",
			})
			return
		}

		helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid image",
		})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext

	avatarDir := "../uploads/avatars"

	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		log.Println(err)

		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not create avatar directory",
		})
		return
	}

	filePath := filepath.Join(avatarDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println(err)

		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not save avatar",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Println(err)

		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not save avatar",
		})
		return
	}

	avatar := filepath.ToSlash(filepath.Join("../uploads/avatars", filename))

	if err := database.UpdateUserAvatar(app.DB, avatar, userId); err != nil {
		log.Println(err)

		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not update avatar",
		})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"status":      true,
		"message":     "avatar updated",
		"avatar_path": avatar,
	})
}

func (app *App) DeleteuserAvatar(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	err := database.DeleteUserAvatar(app.DB, userId)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "could not delete image",
		})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "image deleted",
	})
	return
}
