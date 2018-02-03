package db

import "errors"

var (
	// ErrNotImplemented is the error returned for functionality that is not implemented
	ErrNotImplemented = errors.New("Not implemented")
	// ErrUserExists is the error returned when a user with the same name already exists
	ErrUserExists = errors.New("User already exists")
	// ErrNoSuchUser is the error returned when a user does not exist
	ErrNoSuchUser = errors.New("No user with provided ID")
	// ErrNoSuchFeed is the error returned when a feed does not exist
	ErrNoSuchFeed = errors.New("No feed with provided ID")
	// ErrNotSubscribed is the error returned when a user does not have a feed among the ones they are subscribed to
	ErrNotSubscribed = errors.New("User has no feed with provided ID")
)
