package groups

import (
	"database/sql"
	"errors"
)

var ErrGroupNotFound = errors.New("group not found")

// JoinGroup adds a user to the group's membership table. If the group already
// has a chat, the same user is added there too.
func AcceptInvitation(db *sql.DB, userID int, groupID int64) error {
	if userID <= 0 || groupID <= 0 {
		return ErrGroupNotFound
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Make sure the group exists
	var exists int

	err = tx.QueryRow(`
		SELECT 1
		FROM groups
		WHERE id = ?
	`, groupID).Scan(&exists)

	if err == sql.ErrNoRows {
		return ErrGroupNotFound
	}

	if err != nil {
		return err
	}

	// Accept only an existing pending invitation
	result, err := tx.Exec(`
		UPDATE group_invitations
		SET status = 'accepted'
		WHERE group_id = ?
		  AND user_id = ?
		  AND status = 'pending'
	`, groupID, userID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("pending group invitation not found")
	}

	// Add user as group member
	_, err = tx.Exec(`
		INSERT OR IGNORE INTO group_members (group_id, user_id)
		VALUES (?, ?)
	`, groupID, userID)

	if err != nil {
		return err
	}

	// Find group chat
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
	}

	if err != nil {
		return err
	}

	// Add member to group chat
	_, err = tx.Exec(`
		INSERT OR IGNORE INTO chat_users (user_id, chat_id)
		VALUES (?, ?)
	`, userID, chatID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeclineInvitation(db *sql.DB, userID int, groupID int64) error {
	if userID <= 0 || groupID <= 0 {
		return ErrGroupNotFound
	}

	result, err := db.Exec(`
		UPDATE group_invitations
		SET status = 'declined'
		WHERE group_id = ?
		  AND user_id = ?
		  AND status = 'pending'
	`, groupID, userID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("pending group invitation not found")
	}

	return nil
}

func AcceptJoinRequest(db *sql.DB, creatorID int, requesterID int, groupID int64) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// verify this user is the group creator
	var ownerID int
	err = tx.QueryRow(`
        SELECT creator_id
        FROM groups
        WHERE id = ?
    `, groupID).Scan(&ownerID)

	if err != nil {
		return err
	}

	if ownerID != creatorID {
		return errors.New("only the group creator can accept join requests")
	}

	// verify pending request exists
	result, err := tx.Exec(`
        UPDATE group_join_requests
        SET status = 'accepted'
        WHERE group_id = ?
          AND user_id = ?
          AND status = 'pending'
    `, groupID, requesterID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("pending join request not found")
	}

	// add requester as member
	_, err = tx.Exec(`
        INSERT OR IGNORE INTO group_members (group_id, user_id)
        VALUES (?, ?)
    `, groupID, requesterID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func RejectJoinRequest(db *sql.DB, creatorID int, requesterID int, groupID int64) error {

	var ownerID int

	err := db.QueryRow(`
        SELECT creator_id
        FROM groups
        WHERE id = ?
    `, groupID).Scan(&ownerID)

	if err != nil {
		return err
	}

	if ownerID != creatorID {
		return errors.New("only the group creator can reject join requests")
	}

	result, err := db.Exec(`
        UPDATE group_join_requests
        SET status = 'rejected'
        WHERE group_id = ?
          AND user_id = ?
          AND status = 'pending'
    `, groupID, requesterID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("pending join request not found")
	}

	return nil
}
