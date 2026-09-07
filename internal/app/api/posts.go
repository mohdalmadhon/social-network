package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"social/database/posts"
	"social/internal/helpers"
	"social/internal/models"
)

func (app *App) AddPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "could not read form",
		})
		return
	}

	allowComments, err := strconv.Atoi(r.FormValue("allowComments"))
	if err != nil {
		allowComments = 1
	}

	groupID, err := strconv.Atoi(r.FormValue("groupID"))
	if err != nil {
		groupID = 0
	}

	var taggedPeople []int

	taggedPeopleData := r.FormValue("taggedPeople")

	if taggedPeopleData != "" {
		if err := json.Unmarshal([]byte(taggedPeopleData), &taggedPeople); err != nil {
			helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid tagged people",
			})
			return
		}
	}

	post := models.RegsiterPost{
		UserID:        userID,
		Content:       r.FormValue("content"),
		AllowComments: allowComments,
		GroupID:       groupID,
		Location:      r.FormValue("location"),
		PeopleTagged:  taggedPeople,
	}

	file, header, err := r.FormFile("image")

	if err == nil {
		defer file.Close()

		imagePath, err := helpers.SaveUploads(file, header, "post")

		if err != nil {
			log.Println(err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not save image",
			})
			return
		}

		post.Image_path = imagePath
	}

	if err := posts.GroupExists(app.DB, post.GroupID, userID); err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "group does not exists",
			})
			return
		}
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get group post",
		})
		return
	}
	if err := posts.AddPost(app.DB, post); err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not upload post",
		})
		return
	}

	helpers.WriteJson(w, http.StatusCreated, map[string]any{
		"status":  true,
		"message": "post uploaded successfully",
	})
}
