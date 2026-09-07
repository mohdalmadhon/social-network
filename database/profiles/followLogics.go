package profiles

import (
	"database/sql"
	"fmt"
	"social/internal/models"
)

// SendFollowRequest creates, updates, or deletes a follow relationship.
//
// requestCode:
// -1 = delete request / remove follow
//
//	1 = follow user
//	0 = request to follow private profile
func SendFollowRequest(db *sql.DB, targetID, followerID, requestCode int) error {
	if targetID == followerID {
		return fmt.Errorf("user cannot follow themselves")
	}

	if requestCode != -1 && requestCode != 0 && requestCode != 1 {
		return fmt.Errorf("invalid follow request code: %d", requestCode)
	}

	// -1 always means remove the relationship.
	if requestCode == -1 {
		_, err := db.Exec(`
			DELETE FROM user_followers
			WHERE target_id = ? AND follower_id = ?
		`, targetID, followerID)

		return err
	}

	// Check whether a relationship already exists.
	var status int

	err := db.QueryRow(`
		SELECT status
		FROM user_followers
		WHERE target_id = ? AND follower_id = ?
	`, targetID, followerID).Scan(&status)

	if err == sql.ErrNoRows {
		// No relationship exists, so create one.
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

	// Relationship already exists, so update it.
	_, err = db.Exec(`
		UPDATE user_followers
		SET status = ?
		WHERE target_id = ? AND follower_id = ?
	`, requestCode, targetID, followerID)

	return err
}

func GetFollowers(db *sql.DB, targetID, count int) (map[int]models.UserRegistration, error) {
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
		WHERE uf.target_id = ?
		LIMIT ?
	`, targetID, count)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	followers := make(map[int]models.UserRegistration)

	for rows.Next() {
		var id int
		var follower models.UserRegistration

		err := rows.Scan(
			&id,
			&follower.FirstName,
			&follower.LastName,
			&follower.Avatar,
		)

		if err != nil {
			return nil, err
		}

		followers[id] = follower
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func GetFollowing(db *sql.DB, followerID, count int) (map[int]models.UserRegistration, error) {
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
		WHERE uf.follower_id = ?
		LIMIT ?
	`, followerID, count)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	following := make(map[int]models.UserRegistration)

	for rows.Next() {
		var id int
		var user models.UserRegistration

		err := rows.Scan(
			&id,
			&user.FirstName,
			&user.LastName,
			&user.Avatar,
		)

		if err != nil {
			return nil, err
		}

		following[id] = user
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return following, nil
}
