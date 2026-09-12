package api

import (
	"net/http"
	"path/filepath"
	"social/database/comments"
	"strings"
)

// A private post's image must have the same audience as the post itself.
func (app App) ServeUpload(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	name := strings.TrimPrefix(r.URL.Path, "/uploads/")
	parts := strings.Split(name, "/")
	if len(parts) != 2 || filepath.Base(parts[1]) != parts[1] || strings.Contains(parts[1], "\\") {
		http.NotFound(w, r)
		return
	}
	if parts[0] != "avatars" {
		var postID int64
		var err error
		switch parts[0] {
		case "posts":
			err = app.DB.QueryRow(`SELECT id FROM posts WHERE image_path=?`, name).Scan(&postID)
		case "comments":
			err = app.DB.QueryRow(`SELECT post_id FROM comments WHERE image_path=?`, name).Scan(&postID)
		default:
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.NotFound(w, r)
			return
		}
		allowed, err := comments.CanViewPost(app.DB, userID, postID)
		if err != nil || !allowed {
			http.NotFound(w, r)
			return
		}
	}
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeFile(w, r, filepath.Join("uploads", parts[0], parts[1]))
}
