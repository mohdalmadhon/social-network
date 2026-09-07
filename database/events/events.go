package events

import (
	"database/sql"
	"errors"
)

var ErrEventNotFound = errors.New("event not found")

func SetRSVP(db *sql.DB, userID int, eventID int64, response string) error {
	if userID <= 0 || eventID <= 0 {
		return ErrEventNotFound
	}
	if response != "going" && response != "declined" {
		return errors.New("event response must be going or declined")
	}

	var exists int
	if err := db.QueryRow("SELECT 1 FROM events WHERE id = ?", eventID).Scan(&exists); err == sql.ErrNoRows {
		return ErrEventNotFound
	} else if err != nil {
		return err
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
