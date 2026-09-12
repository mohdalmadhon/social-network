package routes

import (
	"database/sql"
	"net/http"
	"social/internal/app/api"
)

func StartServer(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	app := api.App{DB: db}

	//user
	mux.HandleFunc("GET /api/user", app.AuthMiddleware(app.GetUserData))
	mux.HandleFunc("POST /api/user", app.RegisterUser)
	mux.HandleFunc("PATCH /api/user", app.AuthMiddleware(app.UpdateUserInfo))

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
	mux.HandleFunc("GET /api/search", app.AuthMiddleware(app.Search))
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
	mux.HandleFunc("PATCH /api/events/{eventID}/rsvp", app.AuthMiddleware(app.EventRSVP))
	mux.HandleFunc("POST /api/groups/{id}/invitations", app.AuthMiddleware(app.InviteGroupMember))
	mux.HandleFunc("GET /api/groups/{id}/events", app.AuthMiddleware(app.GroupEvents))
	mux.HandleFunc("POST /api/groups/{id}/events", app.AuthMiddleware(app.GroupEvents))
	mux.HandleFunc("GET /api/groups", app.AuthMiddleware(app.GetGroups))
	mux.HandleFunc("POST /api/groups", app.AuthMiddleware(app.CreateGroup))
	mux.HandleFunc("GET /api/groups/{id}", app.AuthMiddleware(app.GetGroup))
	mux.HandleFunc("DELETE /api/groups/{id}", app.AuthMiddleware(app.DeleteGroup))
	mux.HandleFunc("POST /api/groups/{id}/join-requests", app.AuthMiddleware(app.JoinRequest))
	mux.HandleFunc("DELETE /api/groups/{id}/join-requests", app.AuthMiddleware(app.UndoJoinRequest))

	//folder handlers
	mux.HandleFunc("GET /uploads/", app.AuthMiddleware(app.ServeUpload))
	return mux
}
