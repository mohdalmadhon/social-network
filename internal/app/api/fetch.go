package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	database "social/database/users"
	"social/internal/helpers"
)

func (app *App) SearchLocation(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "no data provided",
		})
		return
	}

	url := "https://nominatim.openstreetmap.org/search" +
		"?format=jsonv2" +
		"&addressdetails=1" +
		"&limit=5" +
		"&q=" + url.QueryEscape(query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)

		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to create location request",
		})
		return
	}

	req.Header.Set("User-Agent", "MySocialNetwork/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to search location",
		})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Nominatim returned status:", resp.StatusCode)

		helpers.WriteJson(w, http.StatusBadGateway, map[string]any{
			"status":  false,
			"message": "location API returned an error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	io.Copy(w, resp.Body)
}

func (app *App) GetFriends(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	searchValue := r.URL.Query().Get("search")
	if searchValue == "" {
		friends, err := database.GetFriends(app.DB, userID)
		if err != nil {
			log.Println(err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not get user friends",
			})
			return
		}
		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status":  true,
			"message": "friends fetched",
			"data":    friends,
		})
		return
	}

	friends, err := database.SearchFriends(app.DB, userID, searchValue)
	if err != nil {
		log.Println(err)
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get friends",
			"data":    friends,
		})
		return
	}
	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "friends fetched",
		"data":    friends,
	})
	return
}
