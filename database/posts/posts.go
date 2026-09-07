package posts

import (
	"database/sql"
	"social/internal/models"
	"strconv"
	"strings"
)

func AddPost(db *sql.DB, post models.RegsiterPost) error {
	tags := func() string {
		tags := make([]string, len(post.PeopleTagged))
		for i, v := range post.PeopleTagged {
			tags[i] = strconv.Itoa(v)
		}
		return strings.Join(tags, ":")
	}

	_, err := db.Exec(`
		INSERT INTO posts (user_id, content, image_path, allow_comments, location, group_id, tags)
		VALUES (?,?,?,?,?,?,?)
	`, post.UserID, post.Content, post.Image_path, post.AllowComments, post.Location, post.GroupID, tags())
	return err
}

func GroupExists(db *sql.DB, groupID, userID int) error {
	if groupID == 0 || groupID == -1 {
		return nil
	}
	err := db.QueryRow(`
		SELECT 1 FROM user_posts_groups WHERE id = ? AND user_id = ?
	`, groupID, userID)
	return err.Err()
}
