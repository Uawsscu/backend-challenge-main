package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("USER_NOT_FOUND")
	ErrEmailAlreadyExists = errors.New("USER_ALREADY_EXISTS")
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")
	ErrInvalidToken       = errors.New("INVALID_TOKEN")
	ErrTokenBlacklisted   = errors.New("TOKEN_BLACKLISTED")
	ErrRequestInvalid     = errors.New("REQUEST_INVALID")
)
