package domain

import (
	"time"
)

// User represents the user entity
type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type TokenClaims struct {
	Subject   string
	ExpiresAt int64
}
