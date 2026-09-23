package users

import (
	"database/sql"
	"social/internal/models"
	"time"
)

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
		VALUES (?,?,?,?,?,NULLIF(?,''))
	`, userData.Email, userData.FirstName, userData.LastName, userData.Password, userData.DOB, userData.UserName)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE profile SET about=?, avatar_path=CASE WHEN ?='' THEN avatar_path ELSE ? END WHERE user_id=?`, userData.About, userData.Avatar, userData.Avatar, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

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
