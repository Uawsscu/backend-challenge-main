package unit

import (
	"context"
	"errors"
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
			password: "Password123!",
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
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

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
		accessToken, refreshToken, err := service.Login(context.Background(), "john@example.com", "Password123!")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if accessToken == "" {
			t.Errorf("expected access token but got empty string")
		}
		if refreshToken == "" {
			t.Errorf("expected refresh token but got empty string")
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
		_, _, err := service.Login(context.Background(), "john@example.com", "wrongpassword")

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
		_, _, err := service.Login(context.Background(), "nonexistent@example.com", "Password123!")

		if err != domain.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials but got %v", err)
		}
	})

	t.Run("token generation error during login", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{ID: "id", Email: email, Password: string(hashedPassword)}, nil
		}

		mockToken.GenerateTokenFunc = func(userID, email string, duration time.Duration) (string, error) {
			return "", errors.New("gen-error")
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.Login(context.Background(), "john@example.com", "Password123!")

		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("session store error during login", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{ID: "id", Email: email, Password: string(hashedPassword)}, nil
		}

		mockToken.GenerateTokenFunc = func(userID, email string, duration time.Duration) (string, error) {
			return "token", nil
		}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "sid"}, nil
		}

		mockSession.StoreSessionFunc = func(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
			return errors.New("redis-error")
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.Login(context.Background(), "john@example.com", "Password123!")

		if err == nil {
			t.Error("expected error but got nil")
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

	t.Run("invalid token", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return nil, domain.ErrInvalidToken
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.ValidateToken(context.Background(), "invalid-token")

		if err != domain.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken but got %v", err)
		}
	})

	t.Run("session not found", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}
		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return "", nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.ValidateToken(context.Background(), "valid-token")

		if err != domain.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken but got %v", err)
		}
	})
}

func TestAuthService_Logout(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockSession := &mocks.MockSessionManager{}
	mockToken := &mocks.MockTokenService{}

	mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
		return &domain.TokenClaims{Subject: "token-id"}, nil
	}

	mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
		return `{"userId":"user-id","refreshToken":"refresh-token-id"}`, nil
	}

	mockSession.DeleteSessionFunc = func(ctx context.Context, key string) error {
		return nil
	}

	service := application.NewAuthService(mockRepo, mockSession, mockToken)
	err := service.Logout(context.Background(), "valid-token")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Run("successful refresh", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `{"userId":"test-user-id","refreshToken":"valid-refresh-token"}`, nil
		}

		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: "test-user-id", Email: "test@example.com"}, nil
		}

		mockToken.GenerateTokenFunc = func(userID, email string, duration time.Duration) (string, error) {
			return "new-token", nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		accessToken, refreshToken, err := service.RefreshToken(context.Background(), "valid-refresh-token")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if accessToken == "" || refreshToken == "" {
			t.Errorf("expected new tokens but got empty")
		}
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return nil, errors.New("invalid")
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.RefreshToken(context.Background(), "invalid-token")

		if err != domain.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken but got %v", err)
		}
	})

	t.Run("unmarshal error", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `invalid-json`, nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.RefreshToken(context.Background(), "valid-token")

		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("token generation error", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `{"userId":"test-user-id","refreshToken":"token"}`, nil
		}

		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: "test-user-id"}, nil
		}

		mockToken.GenerateTokenFunc = func(userID, email string, duration time.Duration) (string, error) {
			return "", errors.New("gen-error")
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.RefreshToken(context.Background(), "token")

		if err == nil || err.Error() != "gen-error" {
			t.Errorf("expected gen-error but got %v", err)
		}
	})
}

func TestAuthService_ValidateToken_EdgeCases(t *testing.T) {
	t.Run("unmarshal error", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `invalid-json`, nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.ValidateToken(context.Background(), "token")

		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("wrong access token", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `{"userId":"id","accessToken":"other-token"}`, nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.ValidateToken(context.Background(), "token")

		if err != domain.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken but got %v", err)
		}
	})

	t.Run("empty user id", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `{"userId":"","accessToken":"token"}`, nil
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		_, _, err := service.ValidateToken(context.Background(), "token")

		if err != domain.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken but got %v", err)
		}
	})
}

func TestAuthService_Logout_EdgeCases(t *testing.T) {
	t.Run("session delete error", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		mockToken.ValidateTokenFunc = func(tokenString string) (*domain.TokenClaims, error) {
			return &domain.TokenClaims{Subject: "token-id"}, nil
		}

		mockSession.GetSessionFunc = func(ctx context.Context, key string) (string, error) {
			return `{"userId":"user-id","refreshToken":"refresh-token-id"}`, nil
		}

		mockSession.DeleteSessionFunc = func(ctx context.Context, key string) error {
			return errors.New("redis-delete-error")
		}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		err := service.Logout(context.Background(), "token")

		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func TestAuthService_Register_EdgeCases(t *testing.T) {
	t.Run("password hash error", func(t *testing.T) {
		mockRepo := &mocks.MockUserRepository{}
		mockSession := &mocks.MockSessionManager{}
		mockToken := &mocks.MockTokenService{}

		service := application.NewAuthService(mockRepo, mockSession, mockToken)
		// Very long password to trigger bcrypt error
		longPassword := string(make([]byte, 80))
		_, err := service.Register(context.Background(), "Name", "email@test.com", longPassword)

		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}
