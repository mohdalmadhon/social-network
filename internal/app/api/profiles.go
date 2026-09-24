package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"social/database/notifications"
	"social/database/profiles"
	"social/database/users"
	"social/internal/helpers"
	"strconv"
)

/*
Handler used to get another user's profile data.

Method:
    GET

-> data should be provided using
 - query parameter: id

-> the id should be the user ID of the profile that needs to be viewed

-> if the profile is public OR the current user is an accepted follower
 - the full profile data will be provided
 - showProfile will be true

-> if the profile is private and the current user is not following
 - limited profile data will be provided
 - showProfile will be false

-> followStatus shows the current follow status
 - 1 = following
 - 0 = follow request sent
 - -1 = not following

-> in case of error there will be a respond written back and can me checked by
 - status boolean
 - message string

-> in case of success a respond will be written back
 - status must be true
 - showProfile boolean
 - followStatus integer
 - data : profile data
*/
func (app *App) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	requestedID := r.URL.Query().Get("id")
	profileID, err := strconv.Atoi(requestedID)
	if err != nil || profileID <= 0 {
		log.Println(err)
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid user id",
		})
		return
	}

	if profileID == userID {
		log.Println("same id")
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "cannot view your own profile this way",
		})
		return
	}

	isPrivate, err := profiles.IsPrivate(app.DB, profileID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJson(w, http.StatusNotFound, map[string]any{
				"status":  false,
				"message": "no user found",
			})
			return
		}

		log.Println("here5", err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not check profile privacy",
		})
		return
	}

	isFollowing, err := profiles.CheckFollower(app.DB, userID, profileID)
	if err != nil {
		if err == sql.ErrNoRows {
			isFollowing = -1
		} else {
			log.Println("here3", err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not check follow status",
			})
			return
		}
	}

	// Public profile OR accepted follower
	if !isPrivate || isFollowing == 1 {
		userData, err := profiles.GetUserData(app.DB, profileID)
		if err != nil {
			if err == sql.ErrNoRows {
				helpers.WriteJson(w, http.StatusNotFound, map[string]any{
					"status":       false,
					"showProfile":  false,
					"followStatus": -1,
					"message":      "no user found",
				})
				return
			}
			log.Println("here1", err)

			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":       false,
				"showProfile":  false,
				"followStatus": -1,
				"message":      "could not get profile data",
			})
			return
		}

		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status":       true,
			"showProfile":  true,
			"followStatus": isFollowing,
			"data":         userData,
		})
		return
	}

	// Private profile and user is not following
	userData, err := profiles.GetPrivateProfileData(app.DB, profileID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJson(w, http.StatusNotFound, map[string]any{
				"status":       false,
				"showProfile":  false,
				"followStatus": -1,
				"message":      "no user found",
			})
			return
		}

		log.Println("here0", err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":       false,
			"showProfile":  false,
			"followStatus": -1,
			"message":      "could not get profile data",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":       true,
		"showProfile":  false,
		"followStatus": isFollowing,
		"data":         userData,
	})
}

/*
Handler used to send a follow request to another user.

Method:
    POST

-> data should be provided using
 - query parameter: targetid

-> if the target user has a private profile
 - a follow request will be created
 - followStatus will be 0

-> if the target user has a public profile
 - the follow will be created directly
 - followStatus will be 1

-> if the target user has a private profile
 - a follow request notification will also be sent

-> in case of error there will be a respond written back and can me checked by
 - status boolean
 - message string

-> in case of success a respond will be written back
 - status must be true
 - followStatus integer
 - message : request sent
*/
func (app *App) RequestFollow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	queryID := r.URL.Query().Get("targetid")
	targetID, err := strconv.Atoi(queryID)
	if err != nil || targetID <= 0 || targetID == followerID {
		log.Println(err)
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid target ID",
		})
		return
	}

	isPrivate, err := profiles.IsPrivate(app.DB, targetID)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to get user data",
		})
		return
	}

	var requestCode int
	if isPrivate {
		requestCode = 0
	} else {
		requestCode = 1
	}
	var previousNotificationID int64
	if requestCode == 0 {
		previousNotification, notificationErr := notifications.GetLatestForActor(
			app.DB,
			targetID,
			"requests",
			"follow_request",
			followerID,
		)
		if notificationErr == nil {
			previousNotificationID = previousNotification.ID
		} else if !errors.Is(notificationErr, sql.ErrNoRows) {
			log.Printf("load previous follow request notification: %v", notificationErr)
		}
	}

	if err := profiles.SendFollowRequest(app.DB, targetID, followerID, requestCode); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "request to follow failed",
		})
		return
	}
	if requestCode == 0 {
		notification, notificationErr := notifications.GetLatestForActor(
			app.DB,
			targetID,
			"requests",
			"follow_request",
			followerID,
		)
		if notificationErr != nil && !errors.Is(notificationErr, sql.ErrNoRows) {
			log.Printf("load follow request notification: %v", notificationErr)
		} else if notificationErr == nil && notification.ID != previousNotificationID {
			app.deliverNotification(notification)
		}
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":       true,
		"followStatus": requestCode,
		"message":      "request sent",
	})
}

/*
Handler used to cancel a follow request.

Method:
    DELETE

-> data should be provided using
 - query parameter: targetid

-> the targetid should be the user ID of the user whose follow request will be removed

-> in case of error there will be a respond written back and can me checked by
 - status boolean
 - message string

-> in case of success a respond will be written back
 - status must be true
 - followStatus will be -1
 - message : request removed
*/
func (app *App) CancelRequest(w http.ResponseWriter, r *http.Request) {
	followerID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	queryID := r.URL.Query().Get("targetid")
	targetID, err := strconv.Atoi(queryID)
	if err != nil || targetID <= 0 || targetID == followerID {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid target ID",
		})
		return
	}

	if err := profiles.SendFollowRequest(app.DB, targetID, followerID, -1); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid target ID",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":       true,
		"followStatus": -1,
		"message":      "request removed",
	})
	return
}

/*
Handler used to remove a follower from the current user's followers.

Method:
    DELETE

-> data should be provided using
 - query parameter: followerid

-> the followerid should be the user ID of the follower that needs to be removed

-> in case of error there will be a respond written back and can me checked by
 - status boolean
 - message string

-> in case of success a respond will be written back
 - status must be true
 - message : follower removed
*/
func (app *App) RemoveFollower(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	followerID, err := strconv.Atoi(r.URL.Query().Get("followerid"))
	if err != nil || followerID <= 0 || followerID == userID {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid follower ID",
		})
		return
	}

	if err := profiles.RemoveFollower(app.DB, userID, followerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteJson(w, http.StatusNotFound, map[string]any{
				"status":  false,
				"message": "follower not found",
			})
			return
		}

		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not remove follower",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "follower removed",
	})
}

/*
Handler used to get the followers of a user.

Method:
    GET

-> data can be provided using
 - query parameter: targetid (optional)
 - query parameter: count (optional)
 - query parameter: offset

-> if targetid is provided
 - the followers of that user will be returned

-> if targetid is not provided
 - the followers of the current logged in user will be returned

-> count is used to control how many followers are returned
 - default value is used if count is not provided
 - the maximum value is limited by maxPageSize

-> offset is used for pagination

-> in case of error there will be a respond written back and can me checked by
 - status boolean
 - message string

-> in case of success a respond will be written back
 - status must be true
 - data : list of followers
*/
func (app *App) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	queryID := r.URL.Query().Get("targetid")
	followerLimit := defaultPageSize
	if rawLimit := r.URL.Query().Get("count"); rawLimit != "" {
		requestedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || requestedLimit < 1 {
			helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid follower count",
			})
			return
		}
		followerLimit = min(requestedLimit, maxPageSize)
	}

	queryOffset := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(queryOffset)

	if err != nil || offset < 0 {
		log.Println(queryOffset)
		log.Println(err, "here1")

		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid offset",
		})
		return
	}

	if queryID != "" && queryID != "null" {
		log.Println(queryID)
		targetID, err := strconv.Atoi(queryID)
		if err != nil || targetID <= 0 {
			log.Println(err, "here2")
			helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid target id",
			})
			return
		}

		followers, err := profiles.GetFollowers(app.DB, targetID, followerLimit, offset)
		if err != nil {
			log.Println(err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not get user followers",
			})
			return
		}

		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status": true,
			"data":   followers,
		})
	} else {
		followers, err := users.GetFollowers(app.DB, userID, followerLimit, offset)
		if err != nil {
			if err == sql.ErrNoRows {
				helpers.WriteJson(w, http.StatusOK, map[string]any{
					"status": false,
					"data":   nil,
				})
				return
			}

			log.Println(err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not get followers",
			})
			return
		}

		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status": true,
			"data":   followers,
		})
		return
	}
}

/*
Handler used to get the users that a user is following.

Method:
    GET

-> data can be provided using
 - query parameter: targetid (optional)
 - query parameter: offset

-> if targetid is provided
 - the following users of that user will be returned

-> if targetid is not provided
 - the following users of the current logged in user will be returned

-> offset is used for pagination

-> the result is limited to 20 users per request

-> in case of error there will be a respond written back and can me checked by
 - status boolean
 - message string

-> in case of success a respond will be written back
 - status must be true
 - data : list of following users
*/
func (app *App) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	queryID := r.URL.Query().Get("targetid")

	queryOffset := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(queryOffset)
	if err != nil || offset < 0 {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid offset",
		})
		return
	}

	if queryID != "" && queryID != "null" {
		targetID, err := strconv.Atoi(queryID)
		if err != nil || targetID <= 0 {
			helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid target id",
			})
			return
		}

		following, err := profiles.GetFollowing(app.DB, targetID, 20, offset)
		if err != nil {
			log.Println(err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not get user followers",
			})
			return
		}

		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status": true,
			"data":   following,
		})
	} else {
		following, err := users.GetFollowing(app.DB, userID, 20, offset)
		if err != nil {
			if err == sql.ErrNoRows {
				helpers.WriteJson(w, http.StatusOK, map[string]any{
					"status": false,
					"data":   nil,
				})
				return
			}

			log.Println(err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not get followers",
			})
			return
		}

		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status": true,
			"data":   following,
		})
		return
	}
}