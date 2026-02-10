package unit

import (
	"context"
	"errors"
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
			password: "Password123!",
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
			errorType:   domain.ErrRequestInvalid,
		},
		{
			name:        "password too short",
			userName:    "John Doe",
			email:       "john@example.com",
			password:    "12345",
			mockSetup:   func(repo *mocks.MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrRequestInvalid,
		},
		{
			name:     "email already exists",
			userName: "John Doe",
			email:    "john@example.com",
			password: "Password123!",
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

func TestUserService_ListUsers(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := application.NewUserService(mockRepo)

	t.Run("successful list", func(t *testing.T) {
		expectedUsers := []*domain.User{
			{ID: "1", Name: "User 1"},
			{ID: "2", Name: "User 2"},
		}

		mockRepo.FindAllFunc = func(ctx context.Context) ([]*domain.User, error) {
			return expectedUsers, nil
		}

		users, err := service.ListUsers(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users but got %d", len(users))
		}
	})
}

func TestUserService_GetUserCount(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := application.NewUserService(mockRepo)

	t.Run("successful count", func(t *testing.T) {
		mockRepo.CountFunc = func(ctx context.Context) (int64, error) {
			return 50, nil
		}

		count, err := service.GetUserCount(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count != 50 {
			t.Errorf("expected count 50 but got %d", count)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.CountFunc = func(ctx context.Context) (int64, error) {
			return 0, errors.New("db error")
		}

		_, err := service.GetUserCount(context.Background())
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func TestUserService_RepositoryErrors(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	service := application.NewUserService(mockRepo)
	ctx := context.Background()
	dbErr := errors.New("db error")

	t.Run("CreateUser repo error", func(t *testing.T) {
		mockRepo.CreateFunc = func(ctx context.Context, user *domain.User) error {
			return dbErr
		}
		_, err := service.CreateUser(ctx, "Name", "test@example.com", "Password123!")
		if err != dbErr {
			t.Errorf("expected %v but got %v", dbErr, err)
		}
	})

	t.Run("GetUserByID repo error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, dbErr
		}
		_, err := service.GetUserByID(ctx, "id")
		if err != dbErr {
			t.Errorf("expected %v but got %v", dbErr, err)
		}
	})

	t.Run("UpdateUser find error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, dbErr
		}
		_, err := service.UpdateUser(ctx, "id", "Name", "email@test.com")
		if err != dbErr {
			t.Errorf("expected %v but got %v", dbErr, err)
		}
	})

	t.Run("UpdateUser update error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, user *domain.User) error {
			return dbErr
		}
		_, err := service.UpdateUser(ctx, "id", "Name", "email@test.com")
		if err != dbErr {
			t.Errorf("expected %v but got %v", dbErr, err)
		}
	})

	t.Run("ListUsers error", func(t *testing.T) {
		mockRepo.FindAllFunc = func(ctx context.Context) ([]*domain.User, error) {
			return nil, dbErr
		}
		_, err := service.ListUsers(ctx)
		if err != dbErr {
			t.Errorf("expected %v but got %v", dbErr, err)
		}
	})
}
