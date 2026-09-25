package posts

import (
	"database/sql"
	"social/database/users"
	"social/internal/models"
)

/*
function to get specefic user posts from the user ID. the function return the posts information and user simple data.
it also needs lmit of how many posts per request and offset.

Parameters:

	db *sql.DB,
	userID int
		-> id of the author id
	offset, limit int

Returns:

	[]models.Posts
		-> in case of error it return nil

	error
		-> nil if success
		-> returns no rows error
*/
func GetUserPosts(db *sql.DB, userID, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	rows, err := db.Query(`
		SELECT p.id, p.content, p.image_path, p.user_id, p.created_at, p.privacy, p.like_count, p.comment_count, p.location,
			   u.id, u.first_name || ' ' || u.last_name AS author, profile.avatar_path
		FROM posts p
		JOIN users AS u ON p.user_id = u.id
		LEFT JOIN profile ON profile.user_id = p.user_id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
		LIMIT ?
		OFFSET ?
	`, userID, limit, offset)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID,
			&p.Content,
			&p.ImagePath,
			&p.UserID,
			&p.CreatedAt,
			&p.Privacy,
			&p.LikeCount,
			&p.CommentCount,
			&p.Location,
			&p.UserID,
			&p.Author,
			&p.AvatarPath,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func FilterPosts(db *sql.DB, posts *[]models.Post, userID int) ([]models.Post, error) {
	var filteredPosts []models.Post
	if len(*posts) == 0 {
		return *posts, nil
	}

	isFollower, err := users.IsFollower(db, userID, (*posts)[0].UserID)
	if err != nil {
		return nil, err
	}

	for _, p := range *posts {
		if p.Privacy == "public" {
			filteredPosts = append(filteredPosts, p)
			continue
		}

		if p.Privacy == "followers" && isFollower {
			filteredPosts = append(filteredPosts, p)
			continue
		}

		if p.Privacy == "selected" {
			var exists int
			err := db.QueryRow(`SELECT 1 FROM post_viewers WHERE post_id = ? AND viewer_id = ?`, p.ID, userID).Scan(&exists)
			if err != nil && err != sql.ErrNoRows {
				continue
			}

			if err == sql.ErrNoRows {
				continue
			}

			filteredPosts = append(filteredPosts, p)
			continue
		}
	}
	return filteredPosts, nil
}
