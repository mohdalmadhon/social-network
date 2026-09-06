package groups

import (
	"database/sql"
	"errors"
)

var ErrGroupNotFound = errors.New("group not found")

// JoinGroup adds a user to the group's membership table. If the group already
// has a chat, the same user is added there too.
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

	if _, err = tx.Exec(`
		INSERT OR IGNORE INTO group_members (group_id, user_id)
		VALUES (?, ?)
	`, groupID, userID); err != nil {
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
		return tx.Commit()
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
