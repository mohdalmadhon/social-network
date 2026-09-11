package notifications

import (
	"database/sql"
	"social/internal/models"
)

func InsertNotification(db *sql.DB, n models.NewNotification) error {
	result, err := db.Exec(`
		INSERT INTO notifications (user_id, message)
		VALUES (?, ?)
	`, n.UserID, n.Message)

	if err != nil {
		return err
	}

	notificationID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO notifications_types (
			notifications_id,
			message_user_id,
			post_id_tag,
			comment_reply_user_id,
			follow_request_user_id,
			follow_request_accept_user_id,
			follow_user_id,
			post_like_user_id,
			post_dislike_user_id,
			comment_like_user_id,
			comment_mention_user_id,
			post_mention_user_id,
			group_invite_user_id,
			group_join_user_id,
			group_accept_user_id,
			event_invite_user_id,
			event_response_user_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		notificationID,
		n.MessageUserID,
		n.PostIDTag,
		n.CommentReplyUserID,
		n.FollowRequestUserID,
		n.FollowRequestAcceptUserID,
		n.FollowUserID,
		n.PostLikeUserID,
		n.PostDislikeUserID,
		n.CommentLikeUserID,
		n.CommentMentionUserID,
		n.PostMentionUserID,
		n.GroupInviteUserID,
		n.GroupJoinUserID,
		n.GroupAcceptUserID,
		n.EventInviteUserID,
		n.EventResponseUserID,
	)

	return err
}
