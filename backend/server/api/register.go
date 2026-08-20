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
	"social/backend/validation"
	"social/sql/database"
	"time"

	"github.com/google/uuid"
)

type App struct {
	DB *sql.DB
}

func (app App) CheckEmailExists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("method not right")
		return
	}

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

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"available": !taken,
	})
}

func (app App) CheckUsernameExists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("wrong method")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("username not avilable")
		return
	}

	taken, err := database.CheckAvilableUsername(app.DB, username)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"usernameAvilable": taken,
	})
}

func (app App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "method not allowed",
		})

		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
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

	err = validation.ValidateUserData(userData)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	HashedPassword, err := authontication.HashPassword(userData.Password)
	if err != nil {
		http.Error(w, "Could not save hash password", http.StatusInternalServerError)
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "error happened getting data",
		})
	}

	userData.Password = HashedPassword
	err = database.InsertUser(app.DB, userData)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("avatar")

	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "Invalid avatar", http.StatusBadRequest)

			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
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
			http.Error(w, "Could not create avatar directory", http.StatusInternalServerError)

			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "error happened getting data",
			})
			return
		}

		filePath := filepath.Join(avatarDir, filename)
		dst, err := os.Create(filePath)
		if err != nil {
			log.Println(err)
			http.Error(w, "Could not save avatar", http.StatusInternalServerError)
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "error happened getting data",
			})
			return
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Could not save avatar", http.StatusInternalServerError)
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "error happened getting data",
			})
			return
		}

		userData.Avatar = filename

		log.Println("Avatar saved:", filePath)

	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "user registered",
	})
}

func (app App) LoggingUser(w http.ResponseWriter, r *http.Request) {
	var userData models.Logger

	if r.Method != http.MethodPost {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
		return
	}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "error happened while fetching data",
		})
		return
	}

	HashedPassword, err := database.GetPasswordByIdentifier(app.DB, userData.Identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "identifier or password is wrong",
			})
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "error happened while authenticating",
		})
		return
	}

	if !authontication.AuthonticateUser(userData, HashedPassword) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "identifier or password is wrong",
		})
		return
	}

	id, err := database.GetUserIDbyIdentifier(app.DB, userData.Identifier)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not get user data",
		})
		return
	}
	token, err := GenerateToken(id, userData.Identifier)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "could not create authentication token",
		})
		return
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(24 * 30 * time.Hour),
	}

	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"status":  true,
		"message": "logged in",
	})
}
