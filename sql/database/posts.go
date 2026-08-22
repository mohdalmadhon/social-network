package database

import (
	"database/sql"
	"errors"
	"fmt"
	"social/backend/models"
)

var (
	ErrSelectedFollowersRequired = errors.New("select at least one follower")
	ErrInvalidPostViewer         = errors.New("selected viewers must follow the post author")
)

func CreatePost(db *sql.DB, userID int, request models.CreatePostRequest) (models.Post, error) {
	tx, err := db.Begin()
	if err != nil {
		return models.Post{}, err
	}
	defer tx.Rollback()

	selectedIDs := uniqueIDs(request.SelectedFollowerIDs)
	if request.Privacy == models.PostPrivacySelected {
		if len(selectedIDs) == 0 {
			return models.Post{}, ErrSelectedFollowersRequired
		}

		for _, viewerID := range selectedIDs {
			var followsAuthor int
			err = tx.QueryRow(`
				SELECT COUNT(*)
				FROM follows
				WHERE follower_id = ? AND following_id = ?
			`, viewerID, userID).Scan(&followsAuthor)
			if err != nil {
				return models.Post{}, err
			}
			if followsAuthor == 0 {
				return models.Post{}, ErrInvalidPostViewer
			}
		}
	}

	result, err := tx.Exec(`
		INSERT INTO posts (type, title, content, image_path, user_id, group_id, privacy)
		VALUES ('post', '', ?, '', ?, NULL, ?)
	`, request.Content, userID, request.Privacy)
	if err != nil {
		return models.Post{}, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return models.Post{}, err
	}

	if request.Privacy == models.PostPrivacySelected {
		for _, viewerID := range selectedIDs {
			_, err = tx.Exec(`
				INSERT INTO post_viewers (post_id, viewer_id)
				VALUES (?, ?)
			`, postID, viewerID)
			if err != nil {
				return models.Post{}, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return models.Post{}, err
	}

	return GetPostByID(db, postID)
}

func GetPostByID(db *sql.DB, postID int64) (models.Post, error) {
	var post models.Post

	err := db.QueryRow(`
		SELECT
			posts.id,
			posts.user_id,
			users.username,
			COALESCE(profile.avatar_path, ''),
			posts.content,
			posts.image_path,
			posts.privacy,
			posts.created_at,
			posts.like_count,
			posts.comment_count
		FROM posts
		JOIN users ON users.id = posts.user_id
		LEFT JOIN profile ON profile.user_id = users.id
		WHERE posts.id = ?
	`, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Author,
		&post.AvatarPath,
		&post.Content,
		&post.ImagePath,
		&post.Privacy,
		&post.CreatedAt,
		&post.LikeCount,
		&post.CommentCount,
	)

	return post, err
}

func ListFeedPosts(db *sql.DB, viewerID int) ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT
			posts.id,
			posts.user_id,
			users.username,
			COALESCE(profile.avatar_path, ''),
			posts.content,
			posts.image_path,
			posts.privacy,
			posts.created_at,
			posts.like_count,
			posts.comment_count
		FROM posts
		JOIN users ON users.id = posts.user_id
		LEFT JOIN profile ON profile.user_id = users.id
		WHERE posts.group_id IS NULL
		AND (
			posts.user_id = ?
			OR posts.privacy = 'public'
			OR (
				posts.privacy = 'followers'
				AND EXISTS (
					SELECT 1
					FROM follows
					WHERE follows.follower_id = ?
					AND follows.following_id = posts.user_id
				)
			)
			OR (
				posts.privacy = 'selected'
				AND EXISTS (
					SELECT 1
					FROM post_viewers
					WHERE post_viewers.post_id = posts.id
					AND post_viewers.viewer_id = ?
				)
			)
		)
		ORDER BY posts.created_at DESC, posts.id DESC
	`, viewerID, viewerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Author,
			&post.AvatarPath,
			&post.Content,
			&post.ImagePath,
			&post.Privacy,
			&post.CreatedAt,
			&post.LikeCount,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func uniqueIDs(ids []int) []int {
	seen := make(map[int]bool)
	unique := make([]int, 0, len(ids))

	for _, id := range ids {
		if id <= 0 || seen[id] {
			continue
		}
		seen[id] = true
		unique = append(unique, id)
	}

	return unique
}

func IsPostPrivacy(value string) bool {
	switch value {
	case models.PostPrivacyPublic, models.PostPrivacyFollowers, models.PostPrivacySelected:
		return true
	default:
		return false
	}
}

func ValidateSelectedIDs(privacy string, ids []int) error {
	if privacy == models.PostPrivacySelected && len(uniqueIDs(ids)) == 0 {
		return ErrSelectedFollowersRequired
	}
	if privacy != models.PostPrivacySelected && len(ids) > 0 {
		return fmt.Errorf("selected followers are only allowed for selected privacy")
	}
	return nil
}
