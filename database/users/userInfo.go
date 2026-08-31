package database

import (
	"database/sql"
	"social/internal/models"
)

// function to get user ID by identifeir
// @returns -1 for not found
// @params identifeir (username or email)
func GetUserID(db *sql.DB, identifier string) int {
	var id int
	err := db.QueryRow(`select id from user where username = ? OR email = ?`, identifier, identifier).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

func GetUserData(db *sql.DB, userID int) (models.UserData, error) {
	var userData models.UserData
	var username sql.NullString
	err := db.QueryRow(`
		SELECT first_name, last_name, email, username, dob FROM user where id = ?
	`, userID).Scan(
		&userData.UserInfo.FirstName,
		&userData.UserInfo.LastName,
		&userData.UserInfo.Email,
		&username,
		&userData.UserInfo.DOB,
	)

	if err != nil {
		return userData, err
	}

	if username.Valid {
		userData.UserInfo.UserName = username.String
	} else {
		userData.UserInfo.UserName = ""
	}

	if err := db.QueryRow(`
		SELECT num_of_followers, num_of_following, num_of_posts, avatar_path, about, is_private from profile where user_id = ?
	`, userID).Scan(
		&userData.NumOfFollowers,
		&userData.NumOfFollowing,
		&userData.NumOfPosts,
		&userData.UserInfo.Avatar,
		&userData.About.Bio,
		&userData.IsPrivate,
	); err != nil {
		return userData, err
	}

	return userData, nil
}

func UpdateUserInfo(db *sql.DB, userID int, userData *models.UserRegistration) error {
    _, err := db.Exec(`
        UPDATE user
        SET first_name = ?, last_name = ?, email = ?, username = ?
        WHERE id = ?
    `, userData.FirstName, userData.LastName, userData.Email, userData.UserName, userID)

	if err != nil {
		return err
	}
	
	_, err = db.Exec(`
		UPDATE profile
		SET about = ?, is_private = ?
		WHERE user_id = ?
	`, userData.About, userData.IsPrivate, userID)

    return err
}

func UpdateUserAvatar(db *sql.DB, userID int, avatar_path string) error {
	_, err := db.Exec(`
		UPDATE profile
		SET avatar_path = ?
		WHERE user_id = ?
	`, avatar_path, userID)
	return err
}