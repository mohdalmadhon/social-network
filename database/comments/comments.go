package comments

import (
	"database/sql"
	"errors"
	"social/internal/models"
	"strings"
)

var ErrPostNotVisible = errors.New("post is not available")

func CreateComment(db *sql.DB, userID int, postID int64, content string) (models.Comment, error) {
	return CreateCommentWithImage(db, userID, postID, content, "")
}

func CreateCommentWithImage(db *sql.DB, userID int, postID int64, content, imagePath string) (models.Comment, error) {
	content = strings.TrimSpace(content)
	if (content == "" && imagePath == "") || len([]rune(content)) > 200 {
		return models.Comment{}, errors.New("comment needs text or an image, with at most 200 characters")
	}

	canView, err := CanViewPost(db, userID, postID)
	if err != nil {
		return models.Comment{}, err
	}
	if !canView {
		return models.Comment{}, ErrPostNotVisible
	}

	result, err := db.Exec(`
		INSERT INTO comments (user_id, post_id, content, image_path)
		VALUES (?, ?, ?, ?)
	`, userID, postID, content, imagePath)
	if err != nil {
		return models.Comment{}, err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return models.Comment{}, err
	}

	return GetCommentByID(db, commentID)
}

func GetCommentByID(db *sql.DB, commentID int64) (models.Comment, error) {
	var comment models.Comment

	err := db.QueryRow(commentQuery+" WHERE comments.id = ?", commentID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Author,
		&comment.AvatarPath,
		&comment.Content,
		&comment.ImagePath,
		&comment.CreatedAt,
	)

	return comment, err
}

func ListComments(db *sql.DB, userID int, postID int64) ([]models.Comment, error) {
	canView, err := CanViewPost(db, userID, postID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, ErrPostNotVisible
	}

	rows, err := db.Query(commentQuery+" WHERE comments.post_id = ? ORDER BY comments.created_at ASC, comments.id ASC", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Author,
			&comment.AvatarPath,
			&comment.Content,
			&comment.ImagePath,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func CanViewPost(db *sql.DB, viewerID int, postID int64) (bool, error) {
	var canView bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM posts
			WHERE posts.id = ?
			AND posts.group_id IS NULL
			AND (
				posts.user_id = ?
				OR posts.privacy = 'public'
				OR (
					posts.privacy = 'followers'
					AND EXISTS (
						SELECT 1
						FROM user_followers
						WHERE follower_id = ?
						AND target_id = posts.user_id
						AND status = 1
					)
				)
				OR (
					posts.privacy = 'selected'
					AND EXISTS (
						SELECT 1
						FROM post_viewers
						WHERE post_id = posts.id
						AND viewer_id = ?
					)
				)
			)
		)
	`, postID, viewerID, viewerID, viewerID).Scan(&canView)

	return canView, err
}

const commentQuery = `
	SELECT
		comments.id,
		comments.post_id,
		comments.user_id,
		COALESCE(users.username, users.first_name || ' ' || users.last_name),
		COALESCE(profile.avatar_path, ''),
		comments.content,
		comments.image_path,
		comments.created_at
	FROM comments
	JOIN user AS users ON users.id = comments.user_id
	LEFT JOIN profile ON profile.user_id = comments.user_id
`
