package events

import (
	"database/sql"
	"errors"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrNotCreator    = errors.New("only the event creator can delete this event")
)

func SetRSVP(db *sql.DB, userID int, eventID int64, response string) error {
	if userID <= 0 || eventID <= 0 {
		return ErrEventNotFound
	}

	var groupID int64
	if err := db.QueryRow(`SELECT e.group_id FROM events e JOIN group_members m ON m.group_id=e.group_id WHERE e.id=? AND m.user_id=?`, eventID, userID).Scan(&groupID); err == sql.ErrNoRows {
		return ErrEventNotFound
	} else if err != nil {
		return err
	}
	return SetGroupRSVP(db, userID, groupID, eventID, response)
}

func SetGroupRSVP(db *sql.DB, userID int, groupID, eventID int64, response string) error {
	if err := requireGroupEvent(db, userID, groupID, eventID); err != nil {
		return err
	}
	if response != "going" && response != "declined" {
		return errors.New("event response must be going or declined")
	}

	_, err := db.Exec(`
		INSERT INTO event_rsvps (event_id, user_id, response)
		VALUES (?, ?, ?)
		ON CONFLICT(event_id, user_id) DO UPDATE SET
			response = excluded.response,
			updated_at = CURRENT_TIMESTAMP
	`, eventID, userID, response)
	return err
}

func RemoveRSVP(db *sql.DB, userID int, groupID, eventID int64) error {
	if err := requireGroupEvent(db, userID, groupID, eventID); err != nil {
		return err
	}
	_, err := db.Exec(`DELETE FROM event_rsvps WHERE event_id = ? AND user_id = ?`, eventID, userID)
	return err
}

func Delete(db *sql.DB, userID int, groupID, eventID int64) error {
	if userID <= 0 || groupID <= 0 || eventID <= 0 {
		return ErrEventNotFound
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var creatorID int
	err = tx.QueryRow(`
		SELECT e.creator_id
		FROM events e
		JOIN group_members gm ON gm.group_id = e.group_id
		WHERE e.id = ? AND e.group_id = ? AND gm.user_id = ?
	`, eventID, groupID, userID).Scan(&creatorID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEventNotFound
	}
	if err != nil {
		return err
	}
	if creatorID != userID {
		return ErrNotCreator
	}

	result, err := tx.Exec(`DELETE FROM events WHERE id = ? AND group_id = ? AND creator_id = ?`, eventID, groupID, userID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrEventNotFound
	}
	return tx.Commit()
}

func requireGroupEvent(db *sql.DB, userID int, groupID, eventID int64) error {
	if userID <= 0 || groupID <= 0 || eventID <= 0 {
		return ErrEventNotFound
	}
	var exists int
	err := db.QueryRow(`
		SELECT 1
		FROM events e
		JOIN group_members gm ON gm.group_id = e.group_id
		WHERE e.id = ? AND e.group_id = ? AND gm.user_id = ?
	`, eventID, groupID, userID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEventNotFound
	}
	return err
}
