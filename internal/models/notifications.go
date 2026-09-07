package models

import "time"

type Notification struct {
	ID        int64     `json:"id"`
	UserID    int       `json:"userId"`
	ActorID   *int      `json:"actorId,omitempty"`
	Category  string    `json:"category"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	RelatedID *int64    `json:"relatedId,omitempty"`
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateNotificationRequest struct {
	ActorID   *int
	Category  string
	Type      string
	Message   string
	RelatedID *int64
}
