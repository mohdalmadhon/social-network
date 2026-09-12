package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social/database/profiles"
	"social/database/users"
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

	userData, err := users.GetUserData(app.DB, userID)
	if err != nil {
		log.Println(err, "here")
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	userAbout, err := profiles.GetUserAbout(app.DB, userID)
	if err != nil {
		log.Println(err, "here1")
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}
	userData.About = userAbout

	followers, err := profiles.GetFollowers(app.DB, userID, 10, 0)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	userData.Followers = followers

	following, err := profiles.GetFollowers(app.DB, userID, 10, 0)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}
	userData.Following = following

	friends, err := users.GetFriends(app.DB, userID, 0)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	userData.Friends = friends
	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status": true,
		"data":   userData,
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

	if err := users.UpdateUserInfo(app.DB, userID, &userData); err != nil {
		log.Println(err)

		status, message := helpers.NormalizeSQLError(err)

		helpers.WriteJson(w, status, map[string]any{
			"status":  false,
			"message": message,
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

	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "avatar must be JPG or PNG",
		})
		return
	}

	avatarPath, err := helpers.SaveUploads(file, header, "avatar")
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not save avatar",
		})
		return
	}

	err = users.UpdateUserAvatar(app.DB, userID, avatarPath)
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

	log.Println(userAbout.Intrests)
	if err := profiles.UpdateUserAbout(app.DB, userID, &userAbout); err != nil {
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

func (app App) GetUsers(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := r.Context().Value("userID").(int)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	rows, err := app.DB.Query(`
		SELECT id, COALESCE(username, ''), first_name, last_name
		FROM user
		WHERE id != ?
		ORDER BY username ASC
	`, currentUserID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not load users",
		})
		return
	}
	defer rows.Close()

	type User struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	users := []User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
		)

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not read users",
			})
			return
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not load users",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": true,
		"users":  users,
	})
}
