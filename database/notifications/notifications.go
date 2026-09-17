package notifications

import (
	"database/sql"
	"errors"
	"social/internal/models"
	"strings"
)

var ErrInvalidCategory = errors.New("notification category must be requests, groups, events, or messages")

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

func List(db *sql.DB, userID int, category string, pagination ...int) ([]models.Notification, error) {
	query := notificationSelect + ` WHERE n.user_id = ?`
	args := []any{userID}

	if category != "" && category != "all" {
		if !IsCategory(category) {
			return nil, ErrInvalidCategory
		}
		query += ` AND n.category = ?`
		args = append(args, category)
	}

	query += ` ORDER BY n.created_at DESC, n.id DESC`
	if len(pagination) >= 2 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, pagination[0], pagination[1])
	} else {
		query += ` LIMIT 100`
	}
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

func UnreadCount(db *sql.DB, userID int, categories ...string) (int, error) {
	category := ""
	if len(categories) > 0 {
		category = categories[0]
	}
	if category != "" && category != "all" && !IsCategory(category) {
		return 0, ErrInvalidCategory
	}

	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = ? AND is_read = 0
	`
	args := []any{userID}
	if category != "" && category != "all" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
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

// MarkMessageNotificationsRead clears message alerts for one private chat.
// Other notification categories and other conversations stay unread.
func MarkMessageNotificationsRead(db *sql.DB, userID int, chatID int64) error {
	_, err := db.Exec(`
		UPDATE notifications
		SET is_read = 1
		WHERE user_id = ?
		  AND category = 'messages'
		  AND type = 'new_message'
		  AND related_id = ?
	`, userID, chatID)
	return err
}

func GetByID(db *sql.DB, userID int, notificationID int64) (models.Notification, error) {
	return scanNotification(
		db.QueryRow(
			notificationSelect+` WHERE n.user_id = ? AND n.id = ?`,
			userID,
			notificationID,
		),
	)
}

func GetLatestForActor(db *sql.DB, userID int, category, notificationType string, actorID int) (models.Notification, error) {
	return scanNotification(
		db.QueryRow(
			notificationSelect+`
				WHERE n.user_id = ?
				  AND n.category = ?
				  AND n.type = ?
				  AND n.actor_id = ?
				ORDER BY n.id DESC
				LIMIT 1
			`,
			userID,
			category,
			notificationType,
			actorID,
		),
	)
}

func ListByRelatedID(db *sql.DB, category, notificationType string, relatedID int64) ([]models.Notification, error) {
	rows, err := db.Query(
		notificationSelect+`
			WHERE n.category = ?
			  AND n.type = ?
			  AND n.related_id = ?
			ORDER BY n.id
		`,
		category,
		notificationType,
		relatedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.Notification, 0)
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, notification)
	}
	return result, rows.Err()
}

func IsCategory(category string) bool {
	switch category {
	case "requests", "groups", "events":
		return true
	case "messages":
		return true
	default:
		return false
	}
}

const notificationSelect = `
	SELECT
		n.id,
		n.user_id,
		n.actor_id,
		n.category,
		n.type,
		n.message,
		n.related_id,
		n.is_read,
		n.created_at,

		gjr.status AS request_status,
		gi.status AS invitation_status,
		(SELECT status FROM user_followers f WHERE n.type = 'follow_request'
		 AND f.target_id = n.user_id AND f.follower_id = n.actor_id) AS follow_status,
		(SELECT response FROM event_rsvps v WHERE n.type = 'event_created'
		 AND v.event_id = n.related_id AND v.user_id = n.user_id) AS event_response

	FROM notifications n

	LEFT JOIN group_join_requests gjr
		ON n.category = 'groups'
		AND n.type = 'join_request'
		AND gjr.id = n.related_id
		AND gjr.user_id = n.actor_id

	LEFT JOIN group_invitations gi
		ON n.category = 'groups'
		AND n.type = 'invitation'
		AND gi.id = n.related_id
		AND gi.user_id = n.user_id
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
		&notification.RequestStatus,
		&notification.InvitationStatus,
		&notification.FollowStatus,
		&notification.EventResponse,
	); err != nil {
		return models.Notification{}, err
	}
	notification.IsRead = isRead == 1
	return notification, nil
}
