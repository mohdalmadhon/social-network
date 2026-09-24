package posts

import "database/sql"

func DeletePost(db *sql.DB, postID, userID int) error {
	_, err := db.Exec(`
		DELETE FROM posts
			WHERE id = ? AND user_id = ?
	`, postID, userID)
	return err
}