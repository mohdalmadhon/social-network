package profiles

import (
	"database/sql"
	"fmt"
	"social/internal/models"
)

/*
SendFollowRequest creates, removes, or maintains a follow relationship
between two users.

Parameters:
	db *sql.DB, targetID int, followerID int, requestCode int

Returns:
	error
	-> nil if the follow request is successfully created or removed
	-> Error if the users are invalid or the database operation fails
*/
func SendFollowRequest(db *sql.DB, targetID, followerID, requestCode int) error {
	if targetID == followerID {
		return fmt.Errorf("user cannot follow themselves")
	}

	if requestCode != -1 && requestCode != 0 && requestCode != 1 {
		return fmt.Errorf("invalid follow request code: %d", requestCode)
	}

	if requestCode == -1 {
		_, err := db.Exec(`
			DELETE FROM user_followers
			WHERE target_id = ? AND follower_id = ?
		`, targetID, followerID)

		return err
	}

	var status sql.NullInt64

	err := db.QueryRow(`
		SELECT status
		FROM user_followers
		WHERE target_id = ? AND follower_id = ?
	`, targetID, followerID).Scan(&status)

	if err == sql.ErrNoRows {
		_, err = db.Exec(`
			INSERT INTO user_followers (
				target_id,
				follower_id,
				status
			)
			VALUES (?, ?, ?)
		`, targetID, followerID, requestCode)

		return err
	}

	if err != nil {
		return err
	}

	/*
	Repeating a follow request does not downgrade an accepted
	relationship or create a duplicate request.
	*/
	return nil
}

/*
DecideFollowRequest accepts or rejects a pending follow request.

Parameters:
	db *sql.DB, targetID int, followerID int, accept bool

Returns:
	error
	-> nil if the pending request is accepted or rejected
	-> sql.ErrNoRows if no pending request exists
	-> Error if the database operation fails
*/
func DecideFollowRequest(db *sql.DB, targetID, followerID int, accept bool) error {
	query := `DELETE FROM user_followers WHERE target_id = ? AND follower_id = ? AND status = 0`

	if accept {
		query = `UPDATE user_followers SET status = 1 WHERE target_id = ? AND follower_id = ? AND status = 0`
	}

	result, err := db.Exec(query, targetID, followerID)

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

/*
RemoveFollower removes an accepted follower relationship.

Parameters:
	db *sql.DB, targetID int, followerID int

Returns:
	error
	-> nil if the follower is successfully removed
	-> sql.ErrNoRows if the accepted relationship does not exist
	-> Error if the users are invalid or the database operation fails
*/
func RemoveFollower(db *sql.DB, targetID, followerID int) error {
	if targetID == followerID {
		return fmt.Errorf("invalid follower")
	}

	result, err := db.Exec(`
		DELETE FROM user_followers
		WHERE target_id = ? AND follower_id = ? AND status = 1
	`, targetID, followerID)

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

/*
GetFollowers retrieves the accepted followers of a user.

Parameters:
	db *sql.DB, targetID int, count int, offset int

Returns:
	map[int]models.UserRegistration
	-> Map containing follower IDs and their basic profile information

	error
	-> nil if successful
	-> Error if the database query fails
*/
func GetFollowers(db *sql.DB, targetID, count, offset int) (map[int]models.UserRegistration, error) {
	rows, err := db.Query(`
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			p.avatar_path
		FROM user_followers uf
		JOIN user u
			ON u.id = uf.follower_id
		LEFT JOIN profile p
			ON p.user_id = u.id
		WHERE uf.target_id = ? AND uf.status = 1
		ORDER BY u.id
		LIMIT ?
		OFFSET ?
	`, targetID, count, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	followers := make(map[int]models.UserRegistration)

	for rows.Next() {
		var id int
		var follower models.UserRegistration
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
			return nil, err
		}

		if firstName.Valid {
			follower.FirstName = firstName.String
		} else {
			follower.FirstName = ""
		}

		if lastName.Valid {
			follower.LastName = lastName.String
		} else {
			follower.LastName = ""
		}

		if avatar.Valid {
			follower.Avatar = avatar.String
		} else {
			follower.Avatar = ""
		}

		followers[id] = follower
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

/*
GetFollowing retrieves the accepted users that a user follows.

Parameters:
	db *sql.DB, followerID int, count int, offset int

Returns:
	map[int]models.UserRegistration
	-> Map containing following user IDs and their basic profile information

	error
	-> nil if successful
	-> Error if the database query fails
*/
func GetFollowing(db *sql.DB, followerID, count, offset int) (map[int]models.UserRegistration, error) {
	rows, err := db.Query(`
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			p.avatar_path
		FROM user_followers uf
		JOIN user u
			ON u.id = uf.target_id
		LEFT JOIN profile p
			ON p.user_id = u.id
		WHERE uf.follower_id = ? AND uf.status = 1
		LIMIT ? OFFSET ?
	`, followerID, count, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	following := make(map[int]models.UserRegistration)

	for rows.Next() {
		var id int
		var user models.UserRegistration
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
			return nil, err
		}

		if firstName.Valid {
			user.FirstName = firstName.String
		} else {
			user.FirstName = ""
		}

		if lastName.Valid {
			user.LastName = lastName.String
		} else {
			user.LastName = ""
		}

		if avatar.Valid {
			user.Avatar = avatar.String
		} else {
			user.Avatar = ""
		}

		following[id] = user
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return following, nil
}

/*
SearchFollows searches for users that the specified user follows.

Parameters:
	db *sql.DB, userID int, searchValue string

Returns:
	[]models.UserRegistration
	-> List of matching users that the user follows

	error
	-> nil if successful
	-> Error if the database query fails
*/
func SearchFollows(db *sql.DB, userID int, searchValue string) ([]models.UserRegistration, error) {
	searchPattern := "%" + searchValue + "%"

	rows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, p.avatar_path
		FROM user u
		LEFT JOIN profile p
			ON p.user_id = u.id
		WHERE EXISTS (
			SELECT 1
			FROM user_followers uf
			WHERE uf.target_id = u.id
			AND uf.follower_id = ?
			AND uf.status = 1
		)
		AND (
			u.first_name LIKE ?
			OR u.last_name LIKE ?
		)
		LIMIT 100
	`, userID, searchPattern, searchPattern)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	follows := make([]models.UserRegistration, 0)

	for rows.Next() {
		var user models.UserRegistration

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Avatar,
		)

		if err != nil {
			return nil, err
		}

		follows = append(follows, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return follows, nil
}

/*
SearchFollowing searches for users who follow the specified user.

Parameters:
	db *sql.DB, userID int, searchValue string

Returns:
	[]models.UserRegistration
	-> List of matching users who follow the specified user

	error
	-> nil if successful
	-> Error if the database query fails
*/
func SearchFollowing(db *sql.DB, userID int, searchValue string) ([]models.UserRegistration, error) {
	searchPattern := "%" + searchValue + "%"

	rows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, p.avatar_path
		FROM user u
		LEFT JOIN profile p
			ON p.user_id = u.id
		WHERE EXISTS (
			SELECT 1
			FROM user_followers uf
			WHERE uf.follower_id = u.id
			AND uf.target_id = ?
			AND uf.status = 1
		)
		AND (
			u.first_name LIKE ?
			OR u.last_name LIKE ?
		)
		LIMIT 100
	`, userID, searchPattern, searchPattern)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	following := make([]models.UserRegistration, 0)

	for rows.Next() {
		var user models.UserRegistration

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Avatar,
		)

		if err != nil {
			return nil, err
		}

		following = append(following, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return following, nil
}