package database

import "database/sql"

// function to get user ID by identifeir
//@returns -1 for not found
//@params identifeir (username or email)
func GetUserID(db *sql.DB, identifier string) int {
	var id int
	err := db.QueryRow(`select id from user where username = ? OR email = ?`, identifier, identifier).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}