package model

import "errors"

var (
	// ErrSessionNotFound Session errors
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionBadRequest  = errors.New("bad request")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound User errors
	ErrUserNotFound = errors.New("user not found")
)
