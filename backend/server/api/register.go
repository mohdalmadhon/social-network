package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"social/backend/authontication"
	"social/backend/models"
	"social/backend/server/api/helpers"
	"social/backend/validation"
	"social/sql/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

type App struct {
	DB *sql.DB
}

func (app *App) CheckEmailExists(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		log.Println("no email provided")
		return
	}

	taken, err := database.CheckAvilableEmail(app.DB, email)
	if err != nil {
		log.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"available": !taken,
	})
}

func (app *App) CheckUsernameExists(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		log.Println("username not available")
		return
	}

	taken, err := database.CheckAvilableUsername(app.DB, username)
	if err != nil {
		log.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"usernameAvailable": !taken,
	})
}

func (app *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "failed to parse information",
		})
		return
	}

	userData := models.RegisterRequest{
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		DOB:       r.FormValue("dob"),
		Username:  r.FormValue("username"),
		About:     r.FormValue("about"),
	}

	userData.Email = strings.ToLower(strings.TrimSpace(userData.Email))
	userData.Username = strings.ToLower(strings.TrimSpace(userData.Username))

	firstName := []rune(strings.TrimSpace(userData.FirstName))
	lastName := []rune(strings.TrimSpace(userData.LastName))

	if len(firstName) > 0 {
		firstName[0] = []rune(strings.ToUpper(string(firstName[0])))[0]

		for i := 1; i < len(firstName); i++ {
			firstName[i] = []rune(strings.ToLower(string(firstName[i])))[0]
		}

		userData.FirstName = string(firstName)
	}

	if len(lastName) > 0 {
		lastName[0] = []rune(strings.ToUpper(string(lastName[0])))[0]

		for i := 1; i < len(lastName); i++ {
			lastName[i] = []rune(strings.ToLower(string(lastName[i])))[0]
		}

		userData.LastName = string(lastName)
	}

	err = validation.ValidateUserData(userData)
	if err != nil {
		log.Println(err)

		helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid user data",
		})
		return
	}

	hashedPassword, err := authontication.HashPassword(userData.Password)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "error happened getting data",
		})
		return
	}

	userData.Password = hashedPassword

	file, header, err := r.FormFile("avatar")
	if err != nil {
		if err != http.ErrMissingFile {
			helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid image",
			})
			return
		}
	} else {
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		filename := uuid.New().String() + ext
		avatarDir := "../uploads/avatars"

		err := os.MkdirAll(avatarDir, 0755)
		if err != nil {
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

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println(err)

			helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not save avatar",
			})
			return
		}

		userData.Avatar = filePath

		log.Println("Avatar saved:", filePath)
		log.Println("Avatar URL path:", userData.Avatar)
	}

	err = database.InsertUser(app.DB, userData)
	if err != nil {
		log.Println(err)

		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not create user",
		})
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, map[string]any{
		"status":  true,
		"message": "user registered",
	})
}

func (app *App) LoggingUser(w http.ResponseWriter, r *http.Request) {
	var userData models.Logger

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "error happened while fetching data",
		})
		return
	}

	hashedPassword, err := database.GetPasswordByIdentifier(app.DB, userData.Identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"status":  false,
				"message": "identifier or password is wrong",
			})
			return
		}

		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "error happened while authenticating",
		})
		return
	}

	if !authontication.AuthonticateUser(userData, hashedPassword) {
		helpers.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "identifier or password is wrong",
		})
		return
	}

	id, err := database.GetUserIDbyIdentifier(app.DB, userData.Identifier)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	token, err := GenerateToken(id, userData.Identifier)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
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

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "logged in",
	})
}

func (app *App) LogOutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			helpers.WriteJSON(w, http.StatusOK, map[string]any{
				"status":  true,
				"message": "already logged out",
			})
			return
		}

		helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid session",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookie.Name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, 
		SameSite: http.SameSiteLaxMode,
	})

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "logged out",
	})
}

func (app App) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}

	var userData models.User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid request body",
		})
		return
	}

	userData.Id = userId

	if userData.Password != "" {
		hashedPass, err := authontication.HashPassword(userData.Password)
		if err != nil {
			log.Println("UpdateProfileAbout error:", err)
			helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not process password",
			})
			return
		}
		userData.Password = hashedPass
	}

	if err := database.UpdateUser(app.DB, userData); err != nil {
		log.Println("UpdateProfileAbout error:", err)
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not update user data",
		})
		return
	}

	if err := database.UpdateProfileAbout(app.DB, userId, userData.About); err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not update profile",
		})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "data updated",
	})
}
