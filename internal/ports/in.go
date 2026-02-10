package ports

import (
	"context"

	"github.com/backend-challenge/user-api/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserCount(ctx context.Context) (int64, error)
}

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, string, error)
}
