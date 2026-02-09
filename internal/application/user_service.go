package application

import (
	"context"
	"fmt"

	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	// Validate input
	if !validator.ValidateRequired(name) {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if !validator.ValidateEmail(email) {
		return nil, fmt.Errorf("%w: invalid email format", domain.ErrInvalidInput)
	}
	if !validator.ValidatePassword(password) {
		return nil, fmt.Errorf("%w: password must be at least 6 characters", domain.ErrInvalidInput)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error) {
	// Validate input
	if !validator.ValidateRequired(name) {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if !validator.ValidateEmail(email) {
		return nil, fmt.Errorf("%w: invalid email format", domain.ErrInvalidInput)
	}

	// Check if user exists
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	user.Name = name
	user.Email = email

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) GetUserCount(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
