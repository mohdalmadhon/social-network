package posts

import (
	"database/sql"
	"social/internal/models"
)

func InsertComment(db *sql.DB, Comment models.Comment) error {
	_, err := db.Exec(`
		INSERT INTO comments (user_id, post_id, content, reply_to)
		VALUES (?,?,?,?)
	`, Comment.User.ID, Comment.PostID, Comment.Content, Comment.RepltTo)
	return err
}

func GetComments(db *sql.DB, postID, limit, offset int, orderBy string) ([]models.Comment, error) {
	var queryOrderBy string

	if orderBy == "latest" {
		queryOrderBy = "c.created_at DESC"
	} else if orderBy == "popular" {
		queryOrderBy = "c.votes DESC"
	} else {
		queryOrderBy = "c.created_at DESC"
	}

	rows, err := db.Query(`
		SELECT 
			u.id,
			u.first_name,
			u.last_name,
			c.id,
			c.content,
			c.post_id,
			c.created_at,
			c.reply_to,
			c.votes
		FROM comments c
		JOIN user u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY `+queryOrderBy+`
		LIMIT ?
		OFFSET ?
	`, postID, limit, offset)

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
			&postID,
			&comment.CreatedAt,
			&comment.RepltTo,
			&comment.Votes,
		)

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
