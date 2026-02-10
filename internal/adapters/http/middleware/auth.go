package middleware

import (
	"net/http"
	"strings"

	"github.com/backend-challenge/user-api/internal/adapters/http/dto"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/internal/ports"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService ports.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "missing_token",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "invalid_token_format",
				Message: "Authorization header must be in format: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]

		_, userID, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			statusCode := http.StatusUnauthorized
			errorType := "invalid_token"

			if err == domain.ErrTokenBlacklisted {
				errorType = "token_blacklisted"
			}

			c.JSON(statusCode, dto.ErrorResponse{
				Error:   errorType,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
