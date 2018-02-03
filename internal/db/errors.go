package db

import "errors"

var (
	ErrUserExists    = errors.New("User already exists")
	ErrNoSuchUser    = errors.New("No user with provided ID")
	ErrNoSuchFeed    = errors.New("No feed with provided ID")
	ErrNotSubscribed = errors.New("User has no feed with provided ID")
)
