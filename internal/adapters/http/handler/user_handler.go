package handler

import (
	"net/http"

	"github.com/backend-challenge/user-api/internal/adapters/http/dto"
	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/backend-challenge/user-api/internal/ports"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"

		if err == domain.ErrEmailAlreadyExists {
			statusCode = http.StatusConflict
			errorType = "email_exists"
		} else if err == domain.ErrInvalidInput {
			statusCode = http.StatusBadRequest
			errorType = "invalid_input"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"

		if err == domain.ErrUserNotFound {
			statusCode = http.StatusNotFound
			errorType = "user_not_found"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	response := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = &dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req.Name, req.Email)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"

		if err == domain.ErrUserNotFound {
			statusCode = http.StatusNotFound
			errorType = "user_not_found"
		} else if err == domain.ErrEmailAlreadyExists {
			statusCode = http.StatusConflict
			errorType = "email_exists"
		} else if err == domain.ErrInvalidInput {
			statusCode = http.StatusBadRequest
			errorType = "invalid_input"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"

		if err == domain.ErrUserNotFound {
			statusCode = http.StatusNotFound
			errorType = "user_not_found"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
