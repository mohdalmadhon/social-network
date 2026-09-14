package models

import (
	"encoding/json"
)

type WSPayload struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type NewNotification struct {
	ID                        int
	Message                   string
	UserID                    int
	MessageUserID             *int
	PostIDTag                 *int
	CommentIDTag              *int
	CommentReplyUserID        *int
	FollowRequestUserID       *int
	FollowRequestAcceptUserID *int
	FollowUserID              *int
	PostLikeUserID            *int
	PostDislikeUserID         *int
	CommentLikeUserID         *int
	CommentMentionUserID      *int
	PostMentionUserID         *int
	GroupInviteUserID         *int
	GroupJoinUserID           *int
	GroupAcceptUserID         *int
	EventInviteUserID         *int
	EventResponseUserID       *int
}

// user ID will be used if groupID is null which means new chat that did not exists before
type IncomingMessage struct {
	Offset     int    `json:"offset"`
	Private    int    `json:"private"`
	UserID     int    `json:"userID"`
	GroupID    int    `json:"groupID"`
	Content    string `json:"content"`
}
