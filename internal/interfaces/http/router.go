package http

import (
	"github.com/backend-challenge/user-api/internal/application"
	"github.com/backend-challenge/user-api/internal/interfaces/http/handler"
	"github.com/backend-challenge/user-api/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userService *application.UserService,
	authService *application.AuthService,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes (authentication required)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(authService))
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}
