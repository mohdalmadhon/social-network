package chats

import (
	"database/sql"
	"social/internal/models"
)

func GetMessages(db *sql.DB, userID, groupID, offset int) ([]models.Message, error) {
	var messages []models.Message

	var exists int

	err := db.QueryRow(`
		SELECT 1
		FROM groups_users
		WHERE group_id = ?
			AND user_id = ?
	`, groupID, userID).Scan(&exists)

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, cp.avatar_path, m.content, m.created_at, m.id
		FROM messages m
		JOIN user u ON u.id = m.sender_id
		JOIN profile p ON p.user_id = u.id
		WHERE m.group_id = ?
		ORDER BY m.created_at ASC
		LIMIT 30 OFFSET ?
	`, groupID, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Message

		err := rows.Scan(
			&m.User.ID,
			&m.User.FirstName,
			&m.User.LastName,
			&m.User.Avatar,
			&m.Message,
			&m.CreatedAt,
			&m.ID,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func ChatExists(db *sql.DB, groupID int) (bool, error) {
	var exists int

	err := db.QueryRow(`
		SELECT 1
		FROM groups
		WHERE id = ?
	`, groupID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func AddMessages(db *sql.DB, content string, userID, groupID int) error {
	_, err := db.Exec(`
		INSERT INTO MESSAGES (content, sender_id, group_id)
		VALUES (?,?,?)
	`, content, userID, groupID)
	return err
} 