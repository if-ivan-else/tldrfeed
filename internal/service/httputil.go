package service

import (
	"net/http"

	"github.com/if-ivan-else/tldrfeed/internal/db"
)

func errorToStatus(e error) int {
	switch e {
	case db.ErrUserExists:
		return http.StatusConflict
	case db.ErrNoSuchFeed:
	case db.ErrNoSuchUser:
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
