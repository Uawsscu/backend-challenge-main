package ports

import (
	"context"
	"time"

	"github.com/backend-challenge/user-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

type SessionManager interface {
	StoreSession(ctx context.Context, key string, data interface{}, ttl time.Duration) error
	GetSession(ctx context.Context, key string) (string, error)
	DeleteSession(ctx context.Context, key string) error
}

type TokenService interface {
	GenerateToken(userID, email string, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (*domain.TokenClaims, error)
}
