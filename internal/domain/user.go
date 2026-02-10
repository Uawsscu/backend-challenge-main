package domain

import (
	"time"
)

// User represents the user entity
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"` // Never expose password in JSON
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}

// AccessTokenSession represents data stored for an access token
type AccessTokenSession struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"` // Link to RT UUID
}

// RefreshTokenSession represents data stored for a refresh token
type RefreshTokenSession struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"` // Link to AT UUID
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
}
