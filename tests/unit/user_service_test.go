package unit

import (
	"context"
	"testing"
	"time"

	"github.com/backend-challenge/user-api/internal/application"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/tests/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		email       string
		password    string
		mockSetup   func(*mocks.MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name:     "successful user creation",
			userName: "John Doe",
			email:    "john@example.com",
			password: "password123",
			mockSetup: func(repo *mocks.MockUserRepository) {
				repo.CreateFunc = func(ctx context.Context, user *domain.User) error {
					user.ID = "test-id-123"
					user.CreatedAt = time.Now()
					return nil
				}
			},
			expectError: false,
		},
		{
			name:        "invalid email format",
			userName:    "John Doe",
			email:       "invalid-email",
			password:    "password123",
			mockSetup:   func(repo *mocks.MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidInput,
		},
		{
			name:        "password too short",
			userName:    "John Doe",
			email:       "john@example.com",
			password:    "12345",
			mockSetup:   func(repo *mocks.MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidInput,
		},
		{
			name:     "email already exists",
			userName: "John Doe",
			email:    "john@example.com",
			password: "password123",
			mockSetup: func(repo *mocks.MockUserRepository) {
				repo.CreateFunc = func(ctx context.Context, user *domain.User) error {
					return domain.ErrEmailAlreadyExists
				}
			},
			expectError: true,
			errorType:   domain.ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockUserRepository{}
			tt.mockSetup(mockRepo)

			service := application.NewUserService(mockRepo)
			user, err := service.CreateUser(context.Background(), tt.userName, tt.email, tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("expected user but got nil")
				}
				if user != nil && user.Email != tt.email {
					t.Errorf("expected email %s but got %s", tt.email, user.Email)
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := application.NewUserService(mockRepo)

	t.Run("user found", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:    "test-id",
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return expectedUser, nil
		}

		user, err := service.GetUserByID(context.Background(), "test-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user.ID != expectedUser.ID {
			t.Errorf("expected user ID %s but got %s", expectedUser.ID, user.ID)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		}

		_, err := service.GetUserByID(context.Background(), "nonexistent-id")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound but got %v", err)
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := application.NewUserService(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		existingUser := &domain.User{
			ID:    "test-id",
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return existingUser, nil
		}

		mockRepo.UpdateFunc = func(ctx context.Context, user *domain.User) error {
			return nil
		}

		user, err := service.UpdateUser(context.Background(), "test-id", "Jane Doe", "jane@example.com")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user.Name != "Jane Doe" {
			t.Errorf("expected name 'Jane Doe' but got %s", user.Name)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		}

		_, err := service.UpdateUser(context.Background(), "nonexistent-id", "Jane Doe", "jane@example.com")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound but got %v", err)
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := application.NewUserService(mockRepo)

	t.Run("successful delete", func(t *testing.T) {
		mockRepo.DeleteFunc = func(ctx context.Context, id string) error {
			return nil
		}

		err := service.DeleteUser(context.Background(), "test-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.DeleteFunc = func(ctx context.Context, id string) error {
			return domain.ErrUserNotFound
		}

		err := service.DeleteUser(context.Background(), "nonexistent-id")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound but got %v", err)
		}
	})
}
