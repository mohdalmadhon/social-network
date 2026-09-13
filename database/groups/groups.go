package groups

import (
	"database/sql"
	"errors"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrInvitationNotFound = errors.New("pending group invitation not found")
)

// AcceptInvitation accepts one specific pending invitation. If the group has
// a chat, the same user is added there too.
func AcceptInvitation(db *sql.DB, userID int, invitationID int64) error {
	if userID <= 0 || invitationID <= 0 {
		return ErrInvitationNotFound
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var groupID int64
	err = tx.QueryRow(`
		SELECT group_id
		FROM group_invitations
		WHERE id = ?
		  AND user_id = ?
		  AND status = 'pending'
	`, invitationID, userID).Scan(&groupID)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrInvitationNotFound
	}

	if err != nil {
		return err
	}

	result, err := tx.Exec(`
		UPDATE group_invitations
		SET status = 'accepted'
		WHERE id = ?
		  AND user_id = ?
		  AND status = 'pending'
	`, invitationID, userID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrInvitationNotFound
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

func DeclineInvitation(db *sql.DB, userID int, invitationID int64) error {
	if userID <= 0 || invitationID <= 0 {
		return ErrInvitationNotFound
	}

	result, err := db.Exec(`
		UPDATE group_invitations
		SET status = 'declined'
		WHERE id = ?
		  AND user_id = ?
		  AND status = 'pending'
	`, invitationID, userID)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrInvitationNotFound
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
