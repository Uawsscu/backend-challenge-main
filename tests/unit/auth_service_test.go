package unit

import (
	"context"
	"testing"
	"time"

	"github.com/backend-challenge/user-api/internal/application"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/tests/mocks"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		email       string
		password    string
		mockSetup   func(*mocks.MockUserRepository)
		expectError bool
	}{
		{
			name:     "successful registration",
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
			name:        "invalid email",
			userName:    "John Doe",
			email:       "invalid-email",
			password:    "password123",
			mockSetup:   func(repo *mocks.MockUserRepository) {},
			expectError: true,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockUserRepository{}
			mockSession := &mocks.MockSessionManager{}
			mockToken := &mocks.MockTokenService{}

			tt.mockSetup(mockRepo)

			service := application.NewAuthService(mockRepo, mockSession, mockToken)
			user, err := service.Register(context.Background(), tt.userName, tt.email, tt.password)

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
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	t.Run("successful login", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:       "test-id",
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: string(hashedPassword),
			}, nil
		}

		mockToken.GenerateTokenFunc = func(userID, email string, duration time.Duration) (string, error) {
			return "test-token", nil
		}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{
				Subject: "token-id",
			}, nil
		}

		mockSession.StoreSessionFunc = func(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
			return nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		accessToken, refreshToken, user, err := service.Login(context.Background(), "john@example.com", "password123")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if accessToken == "" {
			t.Errorf("expected access token but got empty string")
		}
		if refreshToken == "" {
			t.Errorf("expected refresh token but got empty string")
		}
		if user == nil {
			t.Errorf("expected user but got nil")
		}
	})

	t.Run("invalid credentials - wrong password", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:       "test-id",
				Email:    "john@example.com",
				Password: string(hashedPassword),
			}, nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, _, err := service.Login(context.Background(), "john@example.com", "wrongpassword")

		if err != domain.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials but got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, _, err := service.Login(context.Background(), "nonexistent@example.com", "password123")

		if err != domain.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials but got %v", err)
		}
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{
				Subject: "token-id",
			}, nil
		}

		mockSession.IsTokenBlacklistedFunc = func(ctx context.Context, tokenID string) (bool, error) {
			return false, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `{"userId":"test-user-id","accessToken":"valid-token","refreshToken":"refresh-token"}`, nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		claims, userID, err := service.ValidateToken(context.Background(), "valid-token")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if claims == nil {
			t.Errorf("expected claims but got nil")
		}
		if userID != "test-user-id" {
			t.Errorf("expected userID test-user-id but got %s", userID)
		}
	})

	t.Run("blacklisted token", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{
				Subject: "token-id",
			}, nil
		}

		mockSession.IsTokenBlacklistedFunc = func(ctx context.Context, tokenID string) (bool, error) {
			return true, nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.ValidateToken(context.Background(), "blacklisted-token")

		if err != domain.ErrTokenBlacklisted {
			t.Errorf("expected ErrTokenBlacklisted but got %v", err)
		}
	})
}
