package models

import "time"

type ChatList struct {
	ID      int
	Members []UserRegistration
	Name    string
	Avatar  string
}

type PrivateChat struct {
	UserID    int
	GroupID   int
	Avatar    string
	FirstName string
	LastName  string
}

type Message struct {
	ID        int
	User      UserRegistration
	Message   string
	CreatedAt time.Time
}
