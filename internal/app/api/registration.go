package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	database "social/database/users"
	"social/internal/helpers"
	"social/internal/models"
	"social/internal/realtime"
	"social/internal/validation"
	"time"
)

type App struct {
	DB            *sql.DB
	Realtime      *realtime.Hub
	EmailAddress  string
	EmailPassword string
}

/*
Handler used to register a new user.

Method:
    POST

-> data should be provided using
 - multipart/form-data

-> required form data
 - FirstName string
 - LastName string
 - UserName string
 - Email string
 - Password string
 - About string
 - dob string (format: YYYY-MM-DD)
 - VerifyToken string
 - Avatar file (optional)

-> the handler will
 - validate the registration data
 - check if the email was verified using the VerifyToken
 - save the avatar if provided
 - hash the password
 - register the user in the database

-> in case of error there will be a respond written back and can me checked by
 - status boolean (false)
 - message string

-> in case of success a respond will be written back
 - status boolean (true)
 - message string
*/
func (app *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "Bad request, form data is too big",
		})
		return
	}

	dobValue := r.FormValue("dob")

	dob, err := time.Parse("2006-01-02", dobValue)
	if err != nil {
		log.Println(err)

		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "Invalid date of birth",
		})
		return
	}

	userData := models.UserRegistration{
		FirstName: r.FormValue("FirstName"),
		LastName:  r.FormValue("LastName"),
		UserName:  r.FormValue("UserName"),
		Email:     r.FormValue("Email"),
		Password:  r.FormValue("Password"),
		About:     r.FormValue("About"),
		DOB:       dob,
		Avatar:    "",
	}

	err = validation.ValidateRegisterData(&userData)
	if err != nil {
		log.Println(err)

		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	verifyToken := r.FormValue("VerifyToken")

	verified, err := database.HasVerifiedEmail(app.DB, userData.Email, verifyToken, time.Now())
	if err != nil {
		log.Println(err)

		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not check email verification",
		})
		return
	}

	if !verified {
		helpers.WriteJson(w, http.StatusForbidden, map[string]any{
			"status":  false,
			"message": "Please verify your email before registering.",
		})
		return
	}

	file, header, err := r.FormFile("Avatar")
	if err != nil {
		if err != http.ErrMissingFile {
			helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid avatar upload",
			})
			return
		}
	} else {
		defer file.Close()

		avatarPath, err := helpers.SaveUploads(file, header, "avatar")
		if err != nil {
			log.Println(err)

			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not save avatar",
			})
			return
		}

		userData.Avatar = avatarPath
	}

	hashedPassword, err := helpers.HashPassword(userData.Password)
	if err != nil {
		log.Println(err)

		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not hash password",
		})
		return
	}

	userData.Password = hashedPassword

	if err := database.RegisterUser(app.DB, &userData, verifyToken); err != nil {
		log.Println(err)

		if errors.Is(err, database.ErrEmailNotVerified) {
			helpers.WriteJson(w, http.StatusForbidden, map[string]any{
				"status":  false,
				"message": "Please verify your email before registering.",
			})
			return
		}

		status, message := helpers.NormalizeSQLError(err)

		helpers.WriteJson(w, status, map[string]any{
			"status":  false,
			"message": message,
		})
		return
	}

	helpers.WriteJson(w, http.StatusCreated, map[string]any{
		"status":  true,
		"message": "Registration successful",
	})
}
