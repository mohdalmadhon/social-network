package api

import (
	"encoding/json"
	"log"
	"net/http"
	database "social/database/users"
	"social/database/users/profiles"
	"social/internal/helpers"
	"social/internal/models"
	"social/internal/validation"
)

func (app *App) GetUserData(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}
	
	userData, err := database.GetUserData(app.DB, userID)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	userAbout, err := profiles.GetUserAbout(app.DB, userID)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}
	userData.About = userAbout

	followers, err := profiles.GetFollowers(app.DB, userID, 10) 
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	userData.Followers = followers

	following, err := profiles.GetFollowers(app.DB, userID, 10) 
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}
	userData.Following = following

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"data": userData,
	})
}

func (app *App) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	var userData models.UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "could not get data",
		})
		return
	}
	log.Println(userData.IsPrivate)
	if err := validation.ValidateUpdateInfo(&userData); err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid data:" + err.Error(),
		})
		return
	}

	if len(userData.Password) != 0 {
		hashedPassword, err := helpers.HashPassword(userData.Password)
		if err != nil {
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not hash password:" + err.Error(),
			})
			return
		}
		userData.Password = hashedPassword
	}

	if err := database.UpdateUserInfo(app.DB, userID, &userData); err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "Error happened updating data",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "data updated!",
	})
}

func (app *App) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "could not get avatar",
		})
		return
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "avatar must be smaller than 5MB",
		})
		return
	}

	contentType := header.Header.Get("Content-Type")

	if contentType != "image/jpeg" && contentType != "image/png" {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "avatar must be JPG or PNG",
		})
		return
	}

	avatarPath, err := helpers.SaveUploads(file, header)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not save avatar",
		})
		return
	}

	err = database.UpdateUserAvatar(app.DB, userID, avatarPath)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not update avatar",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "avatar updated successfully",
		"avatar":  avatarPath,
	})
}

func (app *App) GetUserAbout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	userProfile, err := profiles.GetUserAbout(app.DB, userID)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user profile" + err.Error(),
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status": true,
		"data":   userProfile,
	})
}

func (app *App) UpdateUserAbout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	var userAbout models.UserAbout
	if err := json.NewDecoder(r.Body).Decode(&userAbout); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "error happened fetching data",
		})
		return
	}

	if err := profiles.UpdateUserAbout(app.DB, userID, userAbout); err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not update user about",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  false,
		"message": "user updated!",
	})
}
