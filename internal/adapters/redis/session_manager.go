package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
	client *redis.Client
}

func NewSessionManager(client *redis.Client) *SessionManager {
	return &SessionManager{
		client: client,
	}
}

func (s *SessionManager) StoreSession(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	fullKey := fmt.Sprintf("session:%s", key)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}
	return s.client.Set(ctx, fullKey, jsonData, ttl).Err()
}

func (s *SessionManager) GetSession(ctx context.Context, key string) (string, error) {
	fullKey := fmt.Sprintf("session:%s", key)
	result, err := s.client.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *SessionManager) DeleteSession(ctx context.Context, key string) error {
	fullKey := fmt.Sprintf("session:%s", key)
	return s.client.Del(ctx, fullKey).Err()
}
