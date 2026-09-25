package users

import (
	"database/sql"
	"social/internal/models"
)

/*
GetFriends retrieves users who mutually follow each other.

Parameters:
	db *sql.DB, userID int, offset int

Returns:
	map[int]models.UserRegistration
	-> Map containing friend IDs and their basic profile information

	error
	-> nil if successful
	-> Error if the database query fails
*/
func GetFriends(db *sql.DB, userID, offset int) (map[int]models.UserRegistration, error) {
	friends := make(map[int]models.UserRegistration)

	rows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, p.avatar_path
		FROM user AS u
		LEFT JOIN profile AS p ON u.id = p.user_id
		WHERE EXISTS (
			SELECT 1
			FROM user_followers AS uf1
			WHERE uf1.follower_id = ?
			  AND uf1.target_id = u.id
			  AND uf1.status = 1
		)
		AND EXISTS (
			SELECT 1
			FROM user_followers AS uf2
			WHERE uf2.follower_id = u.id
			  AND uf2.target_id = ?
			  AND uf2.status = 1
		)
		ORDER BY u.id
		LIMIT 30
		OFFSET ?
	`, userID, userID, offset)

	if err != nil {
		return friends, err
	}

	defer rows.Close()

	for rows.Next() {
		var friend models.UserRegistration
		var id int
		var firstName sql.NullString
		var lastName sql.NullString
		var avatar sql.NullString

		err := rows.Scan(
			&id,
			&firstName,
			&lastName,
			&avatar,
		)

		if err != nil {
			return friends, err
		}

		if firstName.Valid {
			friend.FirstName = firstName.String
		}

		if lastName.Valid {
			friend.LastName = lastName.String
		}

		if avatar.Valid {
			friend.Avatar = avatar.String
		}

		friends[id] = friend
	}

	if err := rows.Err(); err != nil {
		return friends, err
	}

	return friends, nil
}

/*
SearchFriends searches for mutual friends whose first or last name
matches the provided search text.

Parameters:
	db *sql.DB, userID int, search string

Returns:
	map[int]models.UserRegistration
	-> Map containing matching friend IDs and their basic profile information

	error
	-> nil if successful
	-> Error if the database query fails
*/
func SearchFriends(db *sql.DB, userID int, search string) (map[int]models.UserRegistration, error) {
	friends := make(map[int]models.UserRegistration)

	search = "%" + search + "%"

	rows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, p.avatar_path
		FROM user AS u
		LEFT JOIN profile AS p ON u.id = p.user_id
		WHERE EXISTS (
			SELECT 1
			FROM user_followers AS uf1
			WHERE uf1.follower_id = ?
			  AND uf1.target_id = u.id
			  AND uf1.status = 1
		)
		AND EXISTS (
			SELECT 1
			FROM user_followers AS uf2
			WHERE uf2.follower_id = u.id
			  AND uf2.target_id = ?
			  AND uf2.status = 1
		)
		AND (
			u.first_name LIKE ?
			OR u.last_name LIKE ?
		)
		LIMIT 100
	`, userID, userID, search, search)

	if err != nil {
		return friends, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.UserRegistration
		var id int
		var firstName sql.NullString
		var lastName sql.NullString
		var avatar sql.NullString

		err := rows.Scan(
			&id,
			&firstName,
			&lastName,
			&avatar,
		)

		if err != nil {
			return friends, err
		}

		if firstName.Valid {
			user.FirstName = firstName.String
		}

		if lastName.Valid {
			user.LastName = lastName.String
		}

		if avatar.Valid {
			user.Avatar = avatar.String
		}

		friends[id] = user
	}

	if err := rows.Err(); err != nil {
		return friends, err
	}

	return friends, nil
}

/*
GetFollowers retrieves users who follow the specified user.

Parameters:
	db *sql.DB, userID int, limit int, offset int

Returns:
	[]models.UserRegistration
	-> List containing follower information

	error
	-> nil if successful
	-> Error if the database query fails
*/
func GetFollowers(db *sql.DB, userID, limit, offset int) ([]models.UserRegistration, error) {
	var followers []models.UserRegistration

	rows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, p.avatar_path
		FROM user_followers uf
		JOIN user u ON u.id = uf.follower_id
		LEFT JOIN profile p ON p.user_id = u.id
		WHERE uf.target_id = ?
		  AND uf.status = 1
		ORDER BY u.id
		LIMIT ? OFFSET ?
	`, userID, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.UserRegistration
		var firstName sql.NullString
		var lastName sql.NullString
		var avatar sql.NullString

		err := rows.Scan(
			&user.ID,
			&firstName,
			&lastName,
			&avatar,
		)

		if err != nil {
			return nil, err
		}

		if firstName.Valid {
			user.FirstName = firstName.String
		}

		if lastName.Valid {
			user.LastName = lastName.String
		}

		if avatar.Valid {
			user.Avatar = avatar.String
		}

		followers = append(followers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

/*
GetFollowing retrieves users that the specified user follows.

Parameters:
	db *sql.DB, userID int, limit int, offset int

Returns:
	[]models.UserRegistration
	-> List containing following user information

	error
	-> nil if successful
	-> Error if the database query fails
*/
func GetFollowing(db *sql.DB, userID, limit, offset int) ([]models.UserRegistration, error) {
	var following []models.UserRegistration

	rows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, p.avatar_path
		FROM user_followers uf
		JOIN user u ON u.id = uf.target_id
		LEFT JOIN profile p ON p.user_id = u.id
		WHERE uf.follower_id = ?
		  AND uf.status = 1
		ORDER BY u.id
		LIMIT ? OFFSET ?
	`, userID, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.UserRegistration
		var firstName sql.NullString
		var lastName sql.NullString
		var avatar sql.NullString

		err := rows.Scan(
			&user.ID,
			&firstName,
			&lastName,
			&avatar,
		)

		if err != nil {
			return nil, err
		}

		if firstName.Valid {
			user.FirstName = firstName.String
		}

		if lastName.Valid {
			user.LastName = lastName.String
		}

		if avatar.Valid {
			user.Avatar = avatar.String
		}

		following = append(following, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return following, nil
}

/*
function to check if a user (follower) follows another user (target)
NOTE: it returns false in case of error or no follower

Parameters:
	db *sql.DB, followerID, targetID int

Returns:
	bool
		-> false in case of no follow or FAIL
	error
		-> nil if success
*/
func IsFollower(db *sql.DB, followerID, targteID int) (bool, error) {
	var isFollower int

	err := db.QueryRow(`
		SELECT 1
		FROM user_followers
		WHERE follower_id = ? AND target_id = ?
	`, followerID, targteID).Scan(&isFollower)

	if err != nil {
		return false, err
	}

	if isFollower == 1 {
		return true, nil
	}

	return false, nil
}