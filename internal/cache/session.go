package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"full-project-mock/internal/domain/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

type sessionCache struct {
	redis *redis.Client
}

func NewSessionRedisRepository(redis *redis.Client) cache.SessionCache {
	return &sessionCache{redis: redis}
}

func (c *sessionCache) SaveSession(ctx context.Context, s *cache.RefreshSession, ttl time.Duration) error {
	key := buildSessionKey(s.UserID, s.TokenID)

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err = c.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return err
	}

	// Добавим tokenID в индекс (для DeleteAll)
	indexKey := buildIndexKey(s.UserID)
	return c.redis.SAdd(ctx, indexKey, s.TokenID).Err()
}

func (c *sessionCache) GetSession(ctx context.Context, userID int64, tokenID string) (*cache.RefreshSession, error) {
	key := buildSessionKey(userID, tokenID)
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var session cache.RefreshSession
	if err = json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (c *sessionCache) DeleteSession(ctx context.Context, userID int64, tokenID string) error {
	key := buildSessionKey(userID, tokenID)
	indexKey := buildIndexKey(userID)

	pipe := c.redis.TxPipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, indexKey, tokenID)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *sessionCache) DeleteAllUserSessions(ctx context.Context, userID int64) error {
	indexKey := buildIndexKey(userID)

	tokenIDs, err := c.redis.SMembers(ctx, indexKey).Result()
	if err != nil {
		return err
	}

	if len(tokenIDs) == 0 {
		return nil
	}

	pipe := c.redis.TxPipeline()
	for _, tokenID := range tokenIDs {
		key := buildSessionKey(userID, tokenID)
		pipe.Del(ctx, key)
	}
	pipe.Del(ctx, indexKey)

	_, err = pipe.Exec(ctx)
	return err
}

// Формат ключа: auth:refresh:<userID>:<tokenID>
func buildSessionKey(userID int64, tokenID string) string {
	return fmt.Sprintf("auth:refresh:%d:%s", userID, tokenID)
}

// Формат index ключа: auth:refresh:index:<userID>
func buildIndexKey(userID int64) string {
	return fmt.Sprintf("auth:refresh:index:%d", userID)
}
