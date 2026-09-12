package users

import (
	"database/sql"
	"social/internal/models"
)

func RegisterUser(db *sql.DB, userData *models.UserRegistration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
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
