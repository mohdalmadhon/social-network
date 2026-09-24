package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultPageSize = 20
	maxPageSize     = 50
)

var errInvalidPagination = errors.New("pagination values are invalid")

type pageRequest struct {
	Limit  int
	Offset int
}

// parsePage keeps large feeds from being loaded in one response. The client
// can ask for the next page by increasing offset, but cannot bypass the cap.
func parsePage(r *http.Request) (pageRequest, error) {
	page := pageRequest{Limit: defaultPageSize}
	query := r.URL.Query()

	if rawLimit := query.Get("limit"); rawLimit != "" {
		limit, err := strconv.Atoi(rawLimit)
		if err != nil || limit < 1 || limit > maxPageSize {
			return pageRequest{}, errInvalidPagination
		}
		page.Limit = limit
	}

	if rawOffset := query.Get("offset"); rawOffset != "" {
		offset, err := strconv.Atoi(rawOffset)
		if err != nil || offset < 0 {
			return pageRequest{}, errInvalidPagination
		}
		page.Offset = offset
	}

	log.Println(page.Limit, page.Offset)
	return page, nil
}

func trimPage[T any](items []T, page pageRequest, newestFirst bool) ([]T, bool) {
	hasMore := len(items) > page.Limit
	if !hasMore {
		return items, false
	}

	if newestFirst {
		return items[:page.Limit], true
	}

	return items[1:], true
}
