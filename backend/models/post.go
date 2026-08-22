package models

import "time"

const (
	PostPrivacyPublic    = "public"
	PostPrivacyFollowers = "followers"
	PostPrivacySelected  = "selected"
)

type CreatePostRequest struct {
	Content             string `json:"content"`
	Privacy             string `json:"privacy"`
	SelectedFollowerIDs []int  `json:"selectedFollowerIds"`
}

type Post struct {
	ID           int64     `json:"id"`
	UserID       int       `json:"userId"`
	Author       string    `json:"author"`
	AvatarPath   string    `json:"avatarPath"`
	Content      string    `json:"content"`
	ImagePath    string    `json:"imagePath"`
	Privacy      string    `json:"privacy"`
	CreatedAt    time.Time `json:"createdAt"`
	LikeCount    int       `json:"likeCount"`
	CommentCount int       `json:"commentCount"`
}
