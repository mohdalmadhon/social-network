package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"social/backend/models"
)

// inserting user
func InsertUser(db *sql.DB, userData models.RegisterRequest) error {
	_, err := db.Exec(`
		INSERT INTO users (email, first_name, last_name, dob, password, created_at)
		VALUES (?,?,?,?,?,?);
	`, userData.Email, userData.FirstName, userData.LastName, userData.DOB, userData.Password, time.Now())

	if len(userData.Username) != 0 {
		_, err := db.Exec(`insert into users (username) values (?) where id = (select id from users where email = ?)`,
			userData.Username, userData.Email)

		if err != nil {
			return err
		}
	}

	var profile models.Profile
	if len(profile.About) != 0 {
		profile.About = userData.About
	} else {
		profile.About = ""
	}

	if len(profile.Avatar_Path) != 0 {
		profile.Avatar_Path = userData.Avatar
	} else {
		profile.Avatar_Path = ""
	}

	err = InsertProfileData(db, profile, userData.Email)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}

func InsertProfileData(db *sql.DB, profile models.Profile, email string) error {
	_, err := db.Exec(`
		insert into profile (user_id, about, avatar_path, num_of_following, num_of_followers, num_of_posts)
		values ((select id from users where email = ?), ?, ?,0,0,0)
	`, email, profile.About, profile.Avatar_Path)

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
	var username sql.NullString

	err := db.QueryRow(`
		SELECT id, email, username, first_name, last_name, dob, created_at, updated_at
		FROM users
		WHERE username = ? OR email = ?
	`, identifier, identifier).Scan(
		&userData.Id,
		&userData.Email,
		&username,
		&userData.FirstName,
		&userData.LastName,
		&userData.DOB,
		&userData.CreatedAt,
		&userData.Updated_at,
	)

	if err != nil {
		return userData, err
	}

	if username.Valid {
		userData.Username = username.String
	} else {
		userData.Username = ""
	}

	return userData, nil
}

func GetProfileData(db *sql.DB, identifier string) (models.Profile, error) {
	var profile models.Profile
	err := db.QueryRow(`select num_of_followers, num_of_following, num_of_posts, about, avatar_path from profile 
						where user_id = (select id from users where email = ? or username = ?)`, identifier, identifier).Scan(
		&profile.Followers,
		&profile.Following,
		&profile.Posts,
		&profile.About,
		&profile.Avatar_Path,
	)

	if err != nil {
		return profile, err
	}

	return profile, nil
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

// function used by /api/chechEmail endpoint to check if username exists during registration in real time
func CheckAvilableUsername(db *sql.DB, username string) (bool, error) {
	var exists int

	err := db.QueryRow(`select 1 from users where username = ? limit 1`, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// function used to get user id by giving it either email or username
func GetUserIDbyIdentifier(db *sql.DB, identifier string) (int, error) {
	var id int
	err := db.QueryRow(`select id from users where username = ? or email = ?`, identifier, identifier).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}
