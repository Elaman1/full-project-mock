package cache

import (
	"context"
	"fmt"
	"full-project-mock/internal/domain/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

type SessionCache struct {
	redis *redis.Client
}

func NewSessionRedisRepository(redis *redis.Client) cache.SessionCache {
	return &SessionCache{redis: redis}
}

func (session *SessionCache) StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := fmt.Sprintf("auth:refresh:%s", userID)
	return session.redis.Set(ctx, key, token, ttl).Err()
}
