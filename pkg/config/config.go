package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	MongoDBURI         string
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	JWTSecret          string
	JWTAccessTokenSec  int
	JWTRefreshTokenSec int
	ServerPort         string
	LogLevel           string
}

func Load() (*Config, error) {
	jwtAccessSec, err := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_SEC", "900"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_SEC: %w", err)
	}

	jwtRefreshSec, err := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_SEC", "2592000"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TOKEN_SEC: %w", err)
	}

	return &Config{
		MongoDBURI:         getEnv("MONGODB_URI", "mongodb://localhost:27017/userdb"),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		JWTAccessTokenSec:  jwtAccessSec,
		JWTRefreshTokenSec: jwtRefreshSec,
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
