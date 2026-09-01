package database

import (
	"database/sql"
	"social/internal/models"
)

func RegisterUser(db *sql.DB, userData *models.UserRegistration) error {
	_, err := db.Exec(`
		INSERT INTO user (email, first_name, last_name, password, dob)
		VALUES (?,?,?,?,?)
	`, userData.Email, userData.FirstName, userData.LastName, userData.Password, userData.DOB)

	if err != nil {
		return err
	}

	if len(userData.UserName) != 0 {
		_, err := db.Query(`update user set username = ? where id = (select id from user where email = ?)`, userData.UserName, userData.Email)
		if err != nil {
			return err
		}
	}
	
	if userData.Avatar == "" {
		userData.Avatar = "avatar/default.png"
	}

	_, err = db.Exec(`
		UPDATE profile
		SET about = ?, avatar_path = ?
		WHERE user_id = (select id from user where email = ?);
	`, userData.About, userData.Avatar, userData.Email)

	return err
}
