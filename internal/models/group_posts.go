package models

import "time"

type GroupPost struct {
	ID           int64     `json:"id"`
	GroupID      int64     `json:"groupId"`
	UserID       int       `json:"userId"`
	Username     string    `json:"username"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	AvatarPath   string    `json:"avatarPath"`
	Content      string    `json:"content"`
	ImagePath    string    `json:"imagePath"`
	CreatedAt    time.Time `json:"createdAt"`
	CommentCount int       `json:"commentCount"`
	IsOwner      bool      `json:"isOwner"`
}

type GroupPostComment struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"postId"`
	UserID     int       `json:"userId"`
	Username   string    `json:"username"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	AvatarPath string    `json:"avatarPath"`
	Content    string    `json:"content"`
	ImagePath  string    `json:"-"`
	CreatedAt  time.Time `json:"createdAt"`
	IsOwner    bool      `json:"isOwner"`
}
