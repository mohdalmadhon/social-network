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
	Location            string `json:"location"`
	ImagePath           string `json:"-"`
}

type Post struct {
	ID           int64     `json:"id"`
	UserID       int       `json:"userId"`
	Author       string    `json:"author"`
	AvatarPath   string    `json:"avatarPath"`
	Content      string    `json:"content"`
	ImagePath    string    `json:"imagePath"`
	Privacy      string    `json:"privacy"`
	Location     string    `json:"location"`
	CreatedAt    time.Time `json:"createdAt"`
	LikeCount    int       `json:"likeCount"`
	DisLikeCount int       `json:"dislikeCount"`
	Liked        bool      `json:"liked"`
	CommentCount int       `json:"commentCount"`
}
