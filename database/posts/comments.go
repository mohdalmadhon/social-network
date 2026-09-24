package posts

import (
	"database/sql"
	"errors"
)

var ErrCommentNotFound = errors.New("comment not found")

func DeleteComment(db *sql.DB, userID int, postID, commentID int64) error {
	result, err := db.Exec(
		`DELETE FROM comments WHERE id = ? AND post_id = ? AND user_id = ?`,
		commentID, postID, userID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrCommentNotFound
	}

	return nil
}
