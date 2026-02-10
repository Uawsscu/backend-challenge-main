package unit

import (
	"context"
	"testing"
	"time"

	"github.com/backend-challenge/user-api/internal/adapters/redis"
	redisClient "github.com/redis/go-redis/v9"
)

// Note: This test requires a redis server to be running on localhost:6379
// If it fails due to connection, we consider it a skip or environment issue.
func TestSessionManager(t *testing.T) {
	client := redisClient.NewClient(&redisClient.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis is not running on localhost:6379, skipping SessionManager tests")
		return
	}

	mgr := redis.NewSessionManager(client)
	key := "test-key"
	data := map[string]string{"foo": "bar"}
	ttl := 10 * time.Second

	t.Run("Store and Get Session", func(t *testing.T) {
		err := mgr.StoreSession(ctx, key, data, ttl)
		if err != nil {
			t.Errorf("failed to store session: %v", err)
		}

		val, err := mgr.GetSession(ctx, key)
		if err != nil {
			t.Errorf("failed to get session: %v", err)
		}
		if val == "" {
			t.Error("expected session value but got empty")
		}
	})

	t.Run("Delete Session", func(t *testing.T) {
		err := mgr.DeleteSession(ctx, key)
		if err != nil {
			t.Errorf("failed to delete session: %v", err)
		}

		val, err := mgr.GetSession(ctx, key)
		if err != nil {
			t.Errorf("failed to get session after delete: %v", err)
		}
		if val != "" {
			t.Error("expected empty value after delete")
		}
	})
}
