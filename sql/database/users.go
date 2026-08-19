package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"social/backend/models"
)

// inserting user
func InsertUser(db *sql.DB, userData models.User) error {
	_, err := db.Exec(`
		INSERT INTO users (email, username, first_name, last_name, dob, password, created_at)
		VALUES (?,?,?,?,?,?,?);
	`, userData.Email, userData.Username, userData.FirstName, userData.LastName, userData.DOB, userData.Password, time.Now())

	return err
}

// removing user
func RemoveUser(db *sql.DB, userID int) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, userID)
	return err
}

/*
update function user must recieve all the user data wether it is updated or not
becuse it will be too much to created a function or switch case for each
*/
func UpdateUser(db *sql.DB, userData models.User) error {
	result, err := db.Exec(`
		UPDATE users
		SET email = ?, username = ?, first_name = ?, last_name = ?, dob = ?, password = ?
		WHERE id = ?
	`,
		userData.Email,
		userData.Username,
		userData.FirstName,
		userData.LastName,
		userData.DOB,
		userData.Password,
		userData.Id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user with id %d not found", userData.Id)
	}

	return nil
}

// database function that recieve email or username as "identifier" and return user data
func GetUserDataByIdentifier(db *sql.DB, identifier string) (models.User, error) {
	var userData models.User

	err := db.QueryRow(`
		SELECT id, email, username, first_name, last_name, dob, created_at, updated_at
		FROM users
		WHERE username = ? OR email = ?
	`, identifier, identifier).Scan(
		&userData.Id,
		&userData.Email,
		&userData.Username,
		&userData.FirstName,
		&userData.LastName,
		&userData.DOB,
		&userData.CreatedAt,
		&userData.Updated_at,
	)
	if err != nil {
		return userData, err
	}

	return userData, nil
}

func GetPasswordByIdentifier(db *sql.DB, identifier string) (string, error) {
	var password string

	err := db.QueryRow(`SELECT password FROM users WHERE username = ? or email = ?`, identifier, identifier).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

// function used by /api/chechEmail endpoint to check if email exists during registration in real time
func CheckAvilableEmail(db *sql.DB, email string) (bool, error) {
	var exists int

	err := db.QueryRow(`select 1 from users where email = ? limit 1`, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}