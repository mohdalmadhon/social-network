package groups

import (
	"database/sql"
	"errors"
)

var ErrGroupNotFound = errors.New("group not found")

// JoinGroup adds a user to the chat that represents a group.
// The existing schema stores group membership in chat_users, so keeping this
// operation here makes the notification action easy to understand and reuse.
func JoinGroup(db *sql.DB, userID int, groupID int64) error {
	if userID <= 0 || groupID <= 0 {
		return ErrGroupNotFound
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists int
	if err := tx.QueryRow("SELECT 1 FROM groups WHERE id = ?", groupID).Scan(&exists); err == sql.ErrNoRows {
		return ErrGroupNotFound
	} else if err != nil {
		return err
	}

	var chatID int64
	err = tx.QueryRow(`
		SELECT id
		FROM chats
		WHERE group_id = ?
		ORDER BY id
		LIMIT 1
	`, groupID).Scan(&chatID)
	if err == sql.ErrNoRows {
		result, createErr := tx.Exec(`
			INSERT INTO chats (type, group_id)
			VALUES ('group', ?)
		`, groupID)
		if createErr != nil {
			return createErr
		}
		chatID, err = result.LastInsertId()
	} else if err != nil {
		return err
	}

	if _, err = tx.Exec(`
		INSERT OR IGNORE INTO chat_users (user_id, chat_id)
		VALUES (?, ?)
	`, userID, chatID); err != nil {
		return err
	}

	return tx.Commit()
}
