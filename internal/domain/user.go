package domain

import (
	"context"
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

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
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

// SessionManager defines the interface for session storage
type SessionManager interface {
	StoreSession(ctx context.Context, key string, data interface{}, ttl time.Duration) error
	GetSession(ctx context.Context, key string) (string, error)
	DeleteSession(ctx context.Context, key string) error
	BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

// TokenService defines the interface for JWT operations
type TokenService interface {
	GenerateToken(userID, email string, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
}
