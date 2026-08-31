package api

import (
	"database/sql"
	"log"
	"net/http"
	"social/database/users/profiles"
	"social/internal/helpers"
	"strconv"
)

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
	if err != nil || targetID == followerID {
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

	if err := profiles.SendFollowRequest(app.DB, targetID, followerID, requestCode); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "request to follow failed",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":       true,
		"followStatus": requestCode,
		"message":      "request sent",
	})
}

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
	if err != nil || targetID == followerID {
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

func (app *App) GetFollowers(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	queryID := r.URL.Query().Get("targetid")
	queryCount := r.URL.Query().Get("count")

	targetID, err := strconv.Atoi(queryID)
	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid target id",
		})
		return
	}
	count, err := strconv.Atoi(queryCount)
	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid count",
		})
		return
	}

	followers, err := profiles.GetFollowers(app.DB, targetID, count)
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
}

func (app *App) GetFollowing(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	queryID := r.URL.Query().Get("targetid")
	queryCount := r.URL.Query().Get("count")

	targetID, err := strconv.Atoi(queryID)
	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid target id",
		})
		return
	}
	count, err := strconv.Atoi(queryCount)
	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid count",
		})
		return
	}

	following, err := profiles.GetFollowing(app.DB, targetID, count)
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
}
