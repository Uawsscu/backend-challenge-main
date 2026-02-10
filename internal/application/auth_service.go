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

type accessTokenSessionDTO struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshTokenSessionDTO struct {
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
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
		return nil, fmt.Errorf("%w: name is required", domain.ErrRequestInvalid)
	}
	if !validator.ValidateEmail(email) {
		return nil, fmt.Errorf("%w: invalid email format", domain.ErrRequestInvalid)
	}
	if !validator.ValidatePassword(password) {
		return nil, fmt.Errorf("%w: password must be at least 6 characters", domain.ErrRequestInvalid)
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

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return "", "", domain.ErrInvalidCredentials
		}
		return "", "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	// Generate Access Token
	accessToken, err := s.tokenService.GenerateToken(user.ID, user.Email, s.getAccessTokenTTL())
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate Refresh Token
	refreshToken, err := s.tokenService.GenerateToken(user.ID, user.Email, s.getRefreshTokenTTL())
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Extract token IDs to store in Redis
	accessClaims, _ := s.tokenService.ValidateToken(accessToken)
	refreshClaims, _ := s.tokenService.ValidateToken(refreshToken)

	accessSession := &accessTokenSessionDTO{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshClaims.Subject,
	}
	accessKey := fmt.Sprintf("access:%s", accessClaims.Subject)
	if err := s.sessionManager.StoreSession(ctx, accessKey, accessSession, s.getAccessTokenTTL()); err != nil {
		return "", "", fmt.Errorf("failed to store access session: %w", err)
	}

	refreshSession := &refreshTokenSessionDTO{
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}
	refreshKey := fmt.Sprintf("refresh:%s", refreshClaims.Subject)
	if err := s.sessionManager.StoreSession(ctx, refreshKey, refreshSession, s.getRefreshTokenTTL()); err != nil {
		return "", "", fmt.Errorf("failed to store refresh session: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// Validate token
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return err
	}

	accessKey := fmt.Sprintf("access:%s", claims.Subject)
	sessionJSON, _ := s.sessionManager.GetSession(ctx, accessKey)
	if sessionJSON != "" {
		var accessSession accessTokenSessionDTO
		if err := json.Unmarshal([]byte(sessionJSON), &accessSession); err == nil {
			if accessSession.RefreshToken != "" {
				refreshKey := fmt.Sprintf("refresh:%s", accessSession.RefreshToken)
				s.sessionManager.DeleteSession(ctx, refreshKey)
			}
		}
	}

	if err := s.sessionManager.DeleteSession(ctx, accessKey); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, string, error) {
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return nil, "", err
	}

	accessKey := fmt.Sprintf("access:%s", claims.Subject)
	sessionJSON, err := s.sessionManager.GetSession(ctx, accessKey)
	if err != nil {
		return nil, "", err
	}
	if sessionJSON == "" {
		return nil, "", domain.ErrInvalidToken
	}

	var session accessTokenSessionDTO
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal session: %w", err)
	}

	if session.AccessToken != token {
		return nil, "", domain.ErrInvalidToken
	}

	if session.UserID == "" {
		return nil, "", domain.ErrInvalidToken
	}

	return claims, session.UserID, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", domain.ErrInvalidToken
	}

	refreshKey := fmt.Sprintf("refresh:%s", claims.Subject)
	sessionJSON, err := s.sessionManager.GetSession(ctx, refreshKey)
	if err != nil {
		return "", "", err
	}
	if sessionJSON == "" {
		return "", "", domain.ErrInvalidToken
	}

	var refreshSession refreshTokenSessionDTO
	if err := json.Unmarshal([]byte(sessionJSON), &refreshSession); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal refresh session: %w", err)
	}

	if refreshSession.RefreshToken != refreshToken {
		return "", "", domain.ErrInvalidToken
	}
	userID := refreshSession.UserID

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	newAccessToken, err := s.tokenService.GenerateToken(user.ID, user.Email, s.getAccessTokenTTL())
	if err != nil {
		return "", "", err
	}
	latestRefreshToken, err := s.tokenService.GenerateToken(user.ID, user.Email, s.getRefreshTokenTTL())
	if err != nil {
		return "", "", err
	}

	// Cleanup old sessions
	s.sessionManager.DeleteSession(ctx, refreshKey)

	newAccessClaims, _ := s.tokenService.ValidateToken(newAccessToken)
	newRefreshClaims, _ := s.tokenService.ValidateToken(latestRefreshToken)

	newAccessSession := &accessTokenSessionDTO{
		UserID:       user.ID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshClaims.Subject,
	}
	newAccessKey := fmt.Sprintf("access:%s", newAccessClaims.Subject)
	s.sessionManager.StoreSession(ctx, newAccessKey, newAccessSession, s.getAccessTokenTTL())

	newRefreshSession := &refreshTokenSessionDTO{
		RefreshToken: latestRefreshToken,
		UserID:       user.ID,
	}
	newRefreshKey := fmt.Sprintf("refresh:%s", newRefreshClaims.Subject)
	s.sessionManager.StoreSession(ctx, newRefreshKey, newRefreshSession, s.getRefreshTokenTTL())

	return newAccessToken, latestRefreshToken, nil
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
