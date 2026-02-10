package middleware

import (
	"errors"
	"net/http"

	"github.com/backend-challenge/user-api/internal/adapters/http/dto"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			statusCode := http.StatusInternalServerError
			errorType := "internal_error"
			message := err.Error()

			var appErr *domain.AppError
			if errors.As(err, &appErr) {
				errorType = appErr.Status
				switch appErr {
				case domain.ErrUserNotFound:
					statusCode = http.StatusNotFound
				case domain.ErrEmailAlreadyExists:
					statusCode = http.StatusConflict
				case domain.ErrInvalidCredentials:
					statusCode = http.StatusUnauthorized
				case domain.ErrRequestInvalid:
					statusCode = http.StatusBadRequest
				case domain.ErrInvalidToken, domain.ErrTokenBlacklisted:
					statusCode = http.StatusUnauthorized
				}
			}

			c.JSON(statusCode, dto.ErrorResponse{
				Error:   errorType,
				Message: message,
			})
			c.Abort()
		}
	}
}
