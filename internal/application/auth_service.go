package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/backend-challenge/user-api/internal/adapters/jwt"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/internal/ports"
	"github.com/backend-challenge/user-api/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo       ports.UserRepository
	sessionManager ports.SessionManager
	tokenService   ports.TokenService
}

func NewAuthService(
	userRepo ports.UserRepository,
	sessionManager ports.SessionManager,
	tokenService ports.TokenService,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		sessionManager: sessionManager,
		tokenService:   tokenService,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
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

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, *domain.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return "", "", nil, domain.ErrInvalidCredentials
		}
		return "", "", nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, domain.ErrInvalidCredentials
	}

	// Generate Access Token
	accessToken, err := s.tokenService.GenerateToken(user.ID, user.Email, s.getAccessTokenTTL())
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate Refresh Token
	refreshToken, err := s.tokenService.GenerateToken(user.ID, user.Email, s.getRefreshTokenTTL())
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Extract token IDs to store in Redis
	accessClaims, _ := s.tokenService.ValidateToken(accessToken)
	refreshClaims, _ := s.tokenService.ValidateToken(refreshToken)

	// Store Access Token Session
	accessSession := &domain.AccessTokenSession{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshClaims.Subject,
	}
	accessKey := fmt.Sprintf("access:%s", accessClaims.Subject)
	if err := s.sessionManager.StoreSession(ctx, accessKey, accessSession, s.getAccessTokenTTL()); err != nil {
		return "", "", nil, fmt.Errorf("failed to store access session: %w", err)
	}

	// Store Refresh Token Session
	refreshSession := &domain.RefreshTokenSession{
		RefreshToken: refreshToken,
		AccessToken:  accessClaims.Subject,
	}
	refreshKey := fmt.Sprintf("refresh:%s", refreshClaims.Subject)
	if err := s.sessionManager.StoreSession(ctx, refreshKey, refreshSession, s.getRefreshTokenTTL()); err != nil {
		return "", "", nil, fmt.Errorf("failed to store refresh session: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Validate token
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return err
	}

	// Try to get session data to find linked token
	accessKey := fmt.Sprintf("access:%s", claims.Subject)
	sessionJSON, _ := s.sessionManager.GetSession(ctx, accessKey)
	if sessionJSON != "" {
		var accessSession domain.AccessTokenSession
		if err := json.Unmarshal([]byte(sessionJSON), &accessSession); err == nil {
			// Delete linked refresh token session
			if accessSession.RefreshToken != "" {
				refreshKey := fmt.Sprintf("refresh:%s", accessSession.RefreshToken)
				s.sessionManager.DeleteSession(ctx, refreshKey)
			}
		}
	}

	// Delete current session from Redis
	if err := s.sessionManager.DeleteSession(ctx, accessKey); err != nil {
		return err
	}

	// Blacklist current token
	if err := s.sessionManager.BlacklistToken(ctx, claims.Subject, s.getAccessTokenTTL()); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, string, error) {
	// Validate JWT token (opaque check)
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return nil, "", err
	}

	// Check if token is blacklisted in Redis
	blacklisted, err := s.sessionManager.IsTokenBlacklisted(ctx, claims.Subject)
	if err != nil {
		return nil, "", err
	}
	if blacklisted {
		return nil, "", domain.ErrTokenBlacklisted
	}

	// Get corresponding session data from Redis using Subject
	accessKey := fmt.Sprintf("access:%s", claims.Subject)
	sessionJSON, err := s.sessionManager.GetSession(ctx, accessKey)
	if err != nil {
		return nil, "", err
	}
	if sessionJSON == "" {
		return nil, "", domain.ErrInvalidToken
	}

	var session domain.AccessTokenSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Strict Token Verification: Compare the full token string
	if session.AccessToken != token {
		return nil, "", domain.ErrInvalidToken
	}

	if session.UserID == "" {
		return nil, "", domain.ErrInvalidToken
	}

	return claims, session.UserID, nil
}

func (s *AuthService) getAccessTokenTTL() time.Duration {
	if ts, ok := s.tokenService.(*jwt.TokenService); ok {
		return ts.GetAccessTokenDuration()
	}
	return 15 * time.Minute
}

func (s *AuthService) getRefreshTokenTTL() time.Duration {
	if ts, ok := s.tokenService.(*jwt.TokenService); ok {
		return ts.GetRefreshTokenDuration()
	}
	return 30 * 24 * time.Hour
}
