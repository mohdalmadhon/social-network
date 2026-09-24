package helpers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

/*
function used to send a message back to the front end wether success or fail by writting JSON in the header

for the parameter [data] recommended to use { map[string]any } to match the JSON format
Parameters: 
	w http.ResponseWriter, code int, data any

Returns:
	void
*/
func WriteJson(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

/*
this function normalize the sql errors to the user by comparing Repeated patterns

Parameters:
	err error

Returns:
	int 
		-> http status code
	stirng 
		-> normalized error
*/
func NormalizeSQLError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, "The requested information could not be found."
	}

	message := err.Error()

	switch {
	case strings.Contains(message, "UNIQUE constraint failed: user.email"):
		return http.StatusConflict, "This email is already registered."

	case strings.Contains(message, "UNIQUE constraint failed: user.username"):
		return http.StatusConflict, "This username is already taken."

	case strings.Contains(message, "UNIQUE constraint failed: profile.user_id"):
		return http.StatusConflict, "A profile already exists for this user."

	case strings.Contains(message, "FOREIGN KEY constraint failed"):
		return http.StatusBadRequest, "The requested data could not be found."

	case strings.Contains(message, "NOT NULL constraint failed"):
		return http.StatusBadRequest, "Required information is missing."

	case strings.Contains(message, "CHECK constraint failed"):
		return http.StatusBadRequest, "Some of the provided information is invalid."

	default:
		return http.StatusInternalServerError, "Something went wrong. Please try again."
	}
}
