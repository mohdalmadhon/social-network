package users

import (
	"database/sql"
	"social/internal/models"
	"time"
)

/*
RegisterUser registers a new user after verifying their email verification token.

Parameters:
	db *sql.DB, userData *models.UserRegistration, verifyToken string

Returns:
	error
	-> nil if the user is registered successfully
	-> ErrEmailNotVerified if the verification token is missing or invalid
	-> Error if the database transaction or user registration fails
*/
func RegisterUser(db *sql.DB, userData *models.UserRegistration, verifyToken string) error {
	if verifyToken == "" {
		return ErrEmailNotVerified
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	consumed, err := tx.Exec(`
		DELETE FROM email_verifications
		WHERE email = ? AND token = ? AND verified = 1 AND expires_at > ?
	`, userData.Email, verifyToken, time.Now().Unix())

	if err != nil {
		return err
	}

	affected, err := consumed.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrEmailNotVerified
	}

	result, err := tx.Exec(`
		INSERT INTO user (email, first_name, last_name, password, dob, username)
		VALUES (?,?,?,?,?,NULLIF(?, ''))
	`, userData.Email, userData.FirstName, userData.LastName, userData.Password, userData.DOB, userData.UserName)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE profile
		SET about = ?,
		    avatar_path = CASE WHEN ? = '' THEN avatar_path ELSE ? END
		WHERE user_id = ?
	`, userData.About, userData.Avatar, userData.Avatar, id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

/*
CheckUserName checks whether a username already exists in the database.

Parameters:
	db *sql.DB, username string

Returns:
	bool
	-> true if the username already exists
	-> false if the username is available

	error
	-> nil if the check is successful
	-> Error if the database query fails
*/
func CheckUserName(db *sql.DB, username string) (bool, error) {
	var exists int

	err := db.QueryRow(`
		SELECT 1
		FROM user
		WHERE username = ?
		LIMIT 1
	`, username).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

/*
CheckUserEmail checks whether an email address already exists in the database.

Parameters:
	db *sql.DB, username string

Returns:
	bool
	-> true if the email address already exists
	-> false if the email address is available

	error
	-> nil if the check is successful
	-> Error if the database query fails
*/
func CheckUserEmail(db *sql.DB, username string) (bool, error) {
	var exists int

	err := db.QueryRow(`
		SELECT 1
		FROM user
		WHERE email = ?
		LIMIT 1
	`, username).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}