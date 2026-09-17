package api

import (
	"log"

	"social/database/notifications"
	"social/internal/models"
)

func (app *App) deliverNotification(notification models.Notification) {
	if app.Realtime == nil {
		return
	}
	event := realtimeEvent{
		Type:         "notification",
		Notification: &notification,
	}
	if err := app.Realtime.SendToUser(notification.UserID, event); err != nil {
		log.Printf("deliver realtime notification: %v", err)
	}
}

func (app *App) deliverNotificationByID(userID int, notificationID int64) {
	notification, err := notifications.GetByID(app.DB, userID, notificationID)
	if err != nil {
		log.Printf("load realtime notification %d: %v", notificationID, err)
		return
	}
	app.deliverNotification(notification)
}

func (app *App) deliverNotificationsByRelatedID(category, notificationType string, relatedID int64) {
	items, err := notifications.ListByRelatedID(app.DB, category, notificationType, relatedID)
	if err != nil {
		log.Printf("load realtime notifications for %s/%s/%d: %v", category, notificationType, relatedID, err)
		return
	}
	for _, notification := range items {
		app.deliverNotification(notification)
	}
}
