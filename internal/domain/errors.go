package domain

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailAlreadyExists is returned when trying to create a user with an existing email
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidToken is returned when JWT token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenBlacklisted is returned when token has been blacklisted
	ErrTokenBlacklisted = errors.New("token has been blacklisted")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
)
