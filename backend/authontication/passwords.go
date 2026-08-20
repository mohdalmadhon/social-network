package authontication

import (
	"social/backend/models"

	"golang.org/x/crypto/bcrypt"
)

// recieve a password as a string and returns a hashed version with an error
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func AuthonticateUser(userData models.Logger, hashedPassword string) bool {
	if err := (bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userData.Pass))); err != nil {
		return false
	}
	return true
}
