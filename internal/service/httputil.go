package service

import (
	"net/http"

	"github.com/if-ivan-else/tldrfeed/internal/db"
)

func errorToStatus(e error) int {
	switch e {
	case db.ErrNotImplemented:
		return http.StatusNotImplemented

	case db.ErrUserExists:
		return http.StatusConflict

	case db.ErrNoSuchFeed:
		fallthrough
	case db.ErrNoSuchUser:
		fallthrough
	case db.ErrNotSubscribed:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
