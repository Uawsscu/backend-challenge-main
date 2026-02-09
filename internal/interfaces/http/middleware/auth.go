package middleware

import (
	"net/http"
	"strings"

	"github.com/backend-challenge/user-api/internal/application"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *application.AuthService) gin.HandlerFunc {
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

		// Extract token from "Bearer <token>"
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

		// Validate token
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

		// Set user info in context
		c.Set("user_id", userID)

		c.Next()
	}
}
