package models

import "time"

type Comment struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"postId"`
	UserID     int       `json:"userId"`
	Author     string    `json:"author"`
	AvatarPath string    `json:"avatarPath"`
	Content    string    `json:"content"`
	ImagePath  string    `json:"imagePath"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}
