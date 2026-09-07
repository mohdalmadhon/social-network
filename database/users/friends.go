package database

import (
	"database/sql"
	"social/internal/models"
)

func GetFriends(db *sql.DB, userID int) (map[int]models.UserRegistration, error) {
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
        )
        AND EXISTS (
            SELECT 1
            FROM user_followers AS uf2
            WHERE uf2.follower_id = u.id
              AND uf2.target_id = ?
        )
    `, userID, userID)

	if err != nil {
		return friends, err
	}

	defer rows.Close()

	for rows.Next() {
		var friend models.UserRegistration
		var id int
		err := rows.Scan(
			&id,
			&friend.FirstName,
			&friend.LastName,
			&friend.Avatar,
		)

		if err != nil {
			return friends, err
		}

		friends[id] = friend
	}

	if err := rows.Err(); err != nil {
		return friends, err
	}

	return friends, nil
}

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
        )
        AND EXISTS (
            SELECT 1
            FROM user_followers AS uf2
            WHERE uf2.follower_id = u.id
              AND uf2.target_id = ?
        )
        AND (
            u.first_name LIKE ?
            OR u.last_name LIKE ?
        )
    `, userID, userID, search, search)

	if err != nil {
		return friends, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.UserRegistration
		var id int

		err := rows.Scan(
			&id,
			&user.FirstName,
			&user.LastName,
			&user.Avatar,
		)

		if err != nil {
			return friends, err
		}

		friends[id] = user
	}

	if err := rows.Err(); err != nil {
		return friends, err
	}

	return friends, nil
}
