package routes

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"social/internal/app/api"
)

func StartServer(db *sql.DB) *http.ServeMux {
	uploadsDir, err := filepath.Abs("uploads")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	app := api.App{DB: db}

	//user
	mux.HandleFunc("GET /api/user", app.AuthMiddleware(app.GetUserData))
	mux.HandleFunc("POST /api/user", app.RegisterUser)
	mux.HandleFunc("PATCH /api/user", app.AuthMiddleware(app.UpdateUserInfo))
	mux.HandleFunc("GET /api/users", app.AuthMiddleware(app.GetUsers))

	//session
	mux.HandleFunc("POST /api/session", app.LoggingUser)
	mux.HandleFunc("GET /api/session", app.AuthMiddleware(app.AuthorizeSession))
	mux.HandleFunc("DELETE /api/session", app.DeleteSession)

	//profiles
	mux.HandleFunc("PATCH /api/profile/avatar", app.AuthMiddleware(app.UpdateUserAvatar))
	mux.HandleFunc("GET /api/profile/about", app.AuthMiddleware(app.GetUserAbout))
	mux.HandleFunc("PATCH /api/profile/about", app.AuthMiddleware(app.UpdateUserAbout))
	mux.HandleFunc("GET /api/profile", app.AuthMiddleware(app.GetUserProfile))
	mux.HandleFunc("GET /api/friends/", app.AuthMiddleware(app.GetFriends))

	// searches
	mux.HandleFunc("GET /api/profile/follows/search", app.AuthMiddleware(app.SearchFollows))
	mux.HandleFunc("GET /api/profile/following/search", app.AuthMiddleware(app.SearchFollowing))
	mux.HandleFunc("GET /api/location/search", app.SearchLocation)

	// follow handler
	mux.HandleFunc("POST /api/profile/follow", app.AuthMiddleware(app.RequestFollow))
	mux.HandleFunc("DELETE /api/profile/follow", app.AuthMiddleware(app.CancelRequest))
	mux.HandleFunc("GET /api/profile/follow", app.AuthMiddleware(app.GetFollowers))
	mux.HandleFunc("GET /api/profile/following", app.AuthMiddleware(app.GetFollowing))

	//posts
	mux.HandleFunc("/api/posts", app.Posts)
	mux.HandleFunc("/api/posts/{postID}/comments", app.Comments)

	//notifications
	mux.HandleFunc("GET /api/notifications", app.AuthMiddleware(app.Notifications))
	mux.HandleFunc("PATCH /api/notifications/read-all", app.AuthMiddleware(app.MarkAllNotificationsRead))
	mux.HandleFunc("PATCH /api/notifications/{notificationID}/action", app.AuthMiddleware(app.ApplyNotificationAction))
	mux.HandleFunc("PATCH /api/notifications/{notificationID}/read", app.AuthMiddleware(app.MarkNotificationRead))

	//groups
	mux.HandleFunc("GET /api/groups", app.AuthMiddleware(app.GetGroups))
	mux.HandleFunc("POST /api/groups", app.AuthMiddleware(app.CreateGroup))
	mux.HandleFunc("GET /api/groups/{id}", app.AuthMiddleware(app.GetGroup))
	mux.HandleFunc("DELETE /api/groups/{id}", app.AuthMiddleware(app.DeleteGroup))
	mux.HandleFunc("POST /api/groups/{id}/join-requests", app.AuthMiddleware(app.JoinRequest))
	mux.HandleFunc("DELETE /api/groups/{id}/join-requests", app.AuthMiddleware(app.UndoJoinRequest))
	mux.HandleFunc("GET /api/groups/{id}/invite-users", app.AuthMiddleware(app.GetInviteUsers))
	mux.HandleFunc("POST /api/groups/{id}/invitations", app.AuthMiddleware(app.UserInvite))
	mux.HandleFunc("DELETE /api/groups/{id}/invitations/{invitationID}", app.AuthMiddleware(app.UndoInvitation))
	mux.HandleFunc("GET /api/groups/{id}/posts", app.AuthMiddleware(app.GetGroupPosts))
	mux.HandleFunc("POST /api/groups/{id}/posts", app.AuthMiddleware(app.CreateGroupPost))
	mux.HandleFunc("DELETE /api/groups/{id}/posts/{postID}", app.AuthMiddleware(app.DeleteGroupPost))
	mux.HandleFunc("GET /api/groups/{id}/posts/{postID}/comments", app.AuthMiddleware(app.GetGroupPostComments))
	mux.HandleFunc("POST /api/groups/{id}/posts/{postID}/comments", app.AuthMiddleware(app.CreateGroupPostComment))
	mux.HandleFunc("DELETE /api/groups/{id}/posts/{postID}/comments/{commentID}", app.AuthMiddleware(app.DeleteGroupPostComment))

	//folder handlers
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))
	return mux
}
