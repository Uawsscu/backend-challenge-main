package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "github.com/backend-challenge/user-api/internal/adapters/http"
	"github.com/backend-challenge/user-api/internal/adapters/jwt"
	"github.com/backend-challenge/user-api/internal/adapters/mongodb"
	"github.com/backend-challenge/user-api/internal/adapters/redis"
	"github.com/backend-challenge/user-api/internal/application"
	"github.com/backend-challenge/user-api/pkg/config"
	"github.com/backend-challenge/user-api/pkg/logger"
	redisClient "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)
	logger.Info("Starting User Management API", map[string]interface{}{
		"port": cfg.ServerPort,
	})

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to MongoDB
	mongoClient, err := connectMongoDB(ctx, cfg.MongoDBURI)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer mongoClient.Disconnect(ctx)
	logger.Info("Connected to MongoDB successfully")

	// Connect to Redis
	rdb := connectRedis(cfg)
	defer rdb.Close()
	logger.Info("Connected to Redis successfully")

	// Initialize repositories and services
	db := mongoClient.Database("userdb")
	userRepo := mongodb.NewUserRepository(db)
	sessionManager := redis.NewSessionManager(rdb)
	tokenService := jwt.NewTokenService(cfg.JWTSecret, cfg.JWTAccessTokenSec, cfg.JWTRefreshTokenSec)

	userService := application.NewUserService(userRepo)
	authService := application.NewAuthService(userRepo, sessionManager, tokenService)

	// Start background goroutine for user count logging
	go logUserCountPeriodically(ctx, userService)

	// Setup HTTP router
	router := httpHandler.SetupRouter(userService, authService)

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", map[string]interface{}{
			"port": cfg.ServerPort,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", map[string]interface{}{
				"error": err.Error(),
			})
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server gracefully...")

	// Cancel context to stop background goroutines
	cancel()

	// Shutdown server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logger.Info("Server stopped")
}

func connectMongoDB(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

func connectRedis(cfg *config.Config) *redisClient.Client {
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Redis connection failed: %v", err)
	}

	return rdb
}

func logUserCountPeriodically(ctx context.Context, userService *application.UserService) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping user count logging goroutine")
			return
		case <-ticker.C:
			count, err := userService.GetUserCount(ctx)
			if err != nil {
				logger.Error("Failed to get user count", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}
			logger.Info("Total users in database", map[string]interface{}{
				"count": count,
			})
		}
	}
}
