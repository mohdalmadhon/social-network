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
