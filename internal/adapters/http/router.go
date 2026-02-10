package http

import (
	"github.com/backend-challenge/user-api/internal/adapters/http/handler"
	"github.com/backend-challenge/user-api/internal/adapters/http/middleware"
	"github.com/backend-challenge/user-api/internal/ports"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userService ports.UserService,
	authService ports.AuthService,
	lotteryService ports.LotteryService,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandler())

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	lotteryHandler := handler.NewLotteryHandler(lotteryService)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(authService))
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		lotteries := v1.Group("/lotteries")
		lotteries.Use(middleware.AuthMiddleware(authService))
		{
			lotteries.GET("/search", lotteryHandler.Search)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}
