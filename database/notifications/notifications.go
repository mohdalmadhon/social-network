package notifications

import (
	"database/sql"
	"errors"
	"social/internal/models"
	"strings"
)

var ErrInvalidCategory = errors.New("notification category must be requests, groups, or events")

func Create(db *sql.DB, userID int, request models.CreateNotificationRequest) (models.Notification, error) {
	if userID <= 0 || !IsCategory(request.Category) {
		return models.Notification{}, ErrInvalidCategory
	}

	request.Category = strings.TrimSpace(request.Category)
	request.Type = strings.TrimSpace(request.Type)
	request.Message = strings.TrimSpace(request.Message)
	if request.Type == "" || request.Message == "" || len([]rune(request.Message)) > 500 {
		return models.Notification{}, errors.New("notification type and message are required")
	}

	result, err := db.Exec(`
		INSERT INTO notifications (user_id, actor_id, category, type, message, related_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, request.ActorID, request.Category, request.Type, request.Message, request.RelatedID)
	if err != nil {
		return models.Notification{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Notification{}, err
	}

	return GetByID(db, userID, id)
}

func List(db *sql.DB, userID int, category string) ([]models.Notification, error) {
	query := notificationSelect + ` WHERE user_id = ?`
	args := []any{userID}

	if category != "" && category != "all" {
		if !IsCategory(category) {
			return nil, ErrInvalidCategory
		}
		query += ` AND category = ?`
		args = append(args, category)
	}

	query += ` ORDER BY created_at DESC, id DESC LIMIT 100`
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Notification{}
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, notification)
	}

	return result, rows.Err()
}

func UnreadCount(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = ? AND is_read = 0
	`, userID).Scan(&count)
	return count, err
}

func MarkRead(db *sql.DB, userID int, notificationID int64) error {
	_, err := db.Exec(`
		UPDATE notifications
		SET is_read = 1
		WHERE id = ? AND user_id = ?
	`, notificationID, userID)
	return err
}

func MarkAllRead(db *sql.DB, userID int) error {
	_, err := db.Exec(`
		UPDATE notifications
		SET is_read = 1
		WHERE user_id = ? AND is_read = 0
	`, userID)
	return err
}

func GetByID(db *sql.DB, userID int, notificationID int64) (models.Notification, error) {
	return scanNotification(db.QueryRow(notificationSelect+` WHERE user_id = ? AND id = ?`, userID, notificationID))
}

func IsCategory(category string) bool {
	switch category {
	case "requests", "groups", "events":
		return true
	default:
		return false
	}
}

const notificationSelect = `
	SELECT id, user_id, actor_id, category, type, message, related_id, is_read, created_at
	FROM notifications
`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanNotification(row rowScanner) (models.Notification, error) {
	var notification models.Notification
	var isRead int
	if err := row.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.ActorID,
		&notification.Category,
		&notification.Type,
		&notification.Message,
		&notification.RelatedID,
		&isRead,
		&notification.CreatedAt,
	); err != nil {
		return models.Notification{}, err
	}
	notification.IsRead = isRead == 1
	return notification, nil
}
