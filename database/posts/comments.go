package posts

import (
	"database/sql"
	"social/internal/models"
)

func InsertComment(db *sql.DB, comment models.Comment) (models.Comment, error) {
	result, err := db.Exec(`
		INSERT INTO comments (user_id, post_id, content, reply_to)
		VALUES (?, ?, ?, ?)
	`, comment.User.ID, comment.PostID, comment.Content, comment.RepltTo)

	if err != nil {
		return comment, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}

	comment.ID = int(id)

	err = db.QueryRow(`
		SELECT created_at
		FROM comments
		WHERE id = ?
	`, comment.ID).Scan(&comment.CreatedAt)

	if err != nil {
		return comment, err
	}

	err = db.QueryRow(`
		SELECT first_name, last_name
		FROM user
		WHERE id = ?
	`, comment.User.ID).Scan(
		&comment.User.FirstName,
		&comment.User.LastName,
	)

	if err != nil {
		return comment, err
	}

	return comment, nil
}

func GetComments(db *sql.DB, postID, replyTo, limit, offset int, orderBy string) ([]models.Comment, error) {
	queryOrderBy := "c.created_at DESC"

	if orderBy == "popular" {
		queryOrderBy = "c.votes DESC, c.created_at DESC"
	}

	var rows *sql.Rows
	var err error

	if replyTo == 0 {
		rows, err = db.Query(`
			SELECT
				u.id,
				u.first_name,
				u.last_name,
				c.id,
				c.content,
				c.post_id,
				c.created_at,
				c.reply_to,
				c.votes,
				(
					SELECT COUNT(*)
					FROM comments r
					WHERE r.reply_to = c.id
				)
			FROM comments c
			JOIN user u ON u.id = c.user_id
			WHERE c.post_id = ?
			AND c.reply_to IS NULL
			ORDER BY `+queryOrderBy+`
			LIMIT ?
			OFFSET ?
		`, postID, limit, offset)
	} else {
		rows, err = db.Query(`
			SELECT
				u.id,
				u.first_name,
				u.last_name,
				c.id,
				c.content,
				c.post_id,
				c.created_at,
				c.reply_to,
				c.votes,
				(
					SELECT COUNT(*)
					FROM comments r
					WHERE r.reply_to = c.id
				)
			FROM comments c
			JOIN user u ON u.id = c.user_id
			WHERE c.post_id = ?
			AND c.reply_to = ?
			ORDER BY `+queryOrderBy+`
			LIMIT ?
			OFFSET ?
		`, postID, replyTo, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		err := rows.Scan(
			&comment.User.ID,
			&comment.User.FirstName,
			&comment.User.LastName,
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.CreatedAt,
			&comment.RepltTo,
			&comment.Votes,
			&comment.Replies,
		)

		if err != nil {
			return nil, err
		}
		
		err = db.QueryRow(`select avatar_path from profile where user_id = ?`, comment.User.ID).Scan(&comment.User.Avatar)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func DeleteComment(db *sql.DB, commentID, userID int) error {
	result, err := db.Exec(`
		DELETE FROM comments
		WHERE id = ?
		AND user_id = ?
	`, commentID, userID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func VoteComment(db *sql.DB, commentID, userID, vote int) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	var currentVote int

	err = tx.QueryRow(`
		SELECT count
		FROM comment_votes
		WHERE comment_id = ?
		AND user_id = ?
	`, commentID, userID).Scan(&currentVote)

	if err == sql.ErrNoRows {
		_, err = tx.Exec(`
			INSERT INTO comment_votes
			(user_id, comment_id, count)
			VALUES (?, ?, ?)
		`, userID, commentID, vote)

		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE comments
			SET votes = votes + ?
			WHERE id = ?
		`, vote, commentID)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if currentVote == vote {
			_, err = tx.Exec(`
				DELETE FROM comment_votes
				WHERE comment_id = ?
				AND user_id = ?
			`, commentID, userID)

			if err != nil {
				return err
			}

			_, err = tx.Exec(`
				UPDATE comments
				SET votes = votes - ?
				WHERE id = ?
			`, vote, commentID)

			if err != nil {
				return err
			}
		} else {
			_, err = tx.Exec(`
				UPDATE comment_votes
				SET count = ?
				WHERE comment_id = ?
				AND user_id = ?
			`, vote, commentID, userID)

			if err != nil {
				return err
			}

			_, err = tx.Exec(`
				UPDATE comments
				SET votes = votes + ?
				WHERE id = ?
			`, vote*2, commentID)

			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
