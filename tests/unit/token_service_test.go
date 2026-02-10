package unit

import (
	"testing"
	"time"

	"github.com/backend-challenge/user-api/internal/adapters/jwt"
	"github.com/backend-challenge/user-api/internal/domain"
)

func TestTokenService(t *testing.T) {
	secret := "test-secret"
	accessSec := 3600
	refreshSec := 86400
	service := jwt.NewTokenService(secret, accessSec, refreshSec)

	t.Run("Generate and Validate Token", func(t *testing.T) {
		userID := "user-123"
		email := "test@example.com"
		duration := 1 * time.Hour

		token, err := service.GenerateToken(userID, email, duration)
		if err != nil {
			t.Errorf("failed to generate token: %v", err)
		}
		if token == "" {
			t.Error("expected token but got empty string")
		}

		claims, err := service.ValidateToken(token)
		if err != nil {
			t.Errorf("failed to validate token: %v", err)
		}
		if claims.Subject == "" {
			t.Error("expected subject in claims but got empty")
		}
	})

	t.Run("Validate Invalid Token", func(t *testing.T) {
		_, err := service.ValidateToken("invalid.token.here")
		if err != domain.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken but got %v", err)
		}
	})

	t.Run("Durations", func(t *testing.T) {
		if service.GetAccessTokenDuration() != 3600*time.Second {
			t.Errorf("wrong access token duration")
		}
		if service.GetRefreshTokenDuration() != 86400*time.Second {
			t.Errorf("wrong refresh token duration")
		}
	})
}
