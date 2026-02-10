package jwt

import (
	"fmt"
	"time"

	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewTokenService(secretKey string, accessSec, refreshSec int) *TokenService {
	return &TokenService{
		secretKey:            secretKey,
		accessTokenDuration:  time.Duration(accessSec) * time.Second,
		refreshTokenDuration: time.Duration(refreshSec) * time.Second,
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (t *TokenService) GenerateToken(userID, email string, duration time.Duration) (string, error) {
	tokenID := uuid.New().String()
	expirationTime := time.Now().Add(duration)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   tokenID, // Use tokenID as sub
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenService) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	sub, _ := claims.GetSubject()
	return &domain.TokenClaims{
		Subject:   sub,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

func (t *TokenService) GetAccessTokenDuration() time.Duration {
	return t.accessTokenDuration
}

func (t *TokenService) GetRefreshTokenDuration() time.Duration {
	return t.refreshTokenDuration
}
