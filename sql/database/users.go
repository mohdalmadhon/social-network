package queries

import (
	"database/sql"
	"fmt"
	"social/backend/models"
	"time"
)

type App struct {
	DB *sql.DB
}

// inserting user
func (app App) InsertUser(userData models.User) error {
	_, err := app.DB.Exec(`
		INSERT INTO users (email, username, first_name, last_name, dob, password, created_at)
		VALUES (?,?,?,?,?,?,?);
	`, userData.Email, userData.Username, userData.FirstName, userData.LastName, userData.DOB, userData.Password, time.Now())

	return err
}

// removing user
func (app App) RemoveUser(userID int) error {
	_, err := app.DB.Exec(`DELETE FROM users WHERE id = ?`, userID)
	return err
}

/*
update function user must recieve all the user data wether it is updated or not
becuse it will be too much to created a function or switch case for each
*/
func (app App) UpdateUser(userData models.User) error {
	result, err := app.DB.Exec(`
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
func (app App) GetUserDataByIdentifier(identifier string) (models.User, error) {
	var userData models.User

	err := app.DB.QueryRow(`
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

func (app App) GetPasswordByIdentifier(identifier string) (string, error) {
	var password string
	err := app.DB.QueryRow(`SELECT password FROM users WHERE username = ? or email = ?`, identifier, identifier).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}
