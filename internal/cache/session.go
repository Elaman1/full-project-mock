package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"full-project-mock/internal/domain/cache"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type sessionCache struct {
	redis *redis.Client
}

func NewSessionRedisRepository(redis *redis.Client) cache.SessionCache {
	return &sessionCache{redis: redis}
}

func (c *sessionCache) SaveSession(ctx context.Context, s *cache.RefreshSession, ttl time.Duration) error {
	key := buildSessionKey(s.TokenID)

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err = c.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return err
	}

	// Добавим tokenID в индекс (для DeleteAll)
	indexKey := buildIndexKey(s.UserID)
	compound := fmt.Sprintf("%s:%s", s.TokenID, s.TokenHash)
	return c.redis.SAdd(ctx, indexKey, compound).Err()
}

func (c *sessionCache) GetSession(ctx context.Context, tokenID string) (*cache.RefreshSession, error) {
	key := buildSessionKey(tokenID)
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
	key := buildSessionKey(tokenID)
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
	for _, compound := range tokenIDs {
		parts := strings.SplitN(compound, ":", 2)
		tokenID := parts[0]
		tokenHash := parts[1]

		key := buildSessionKey(tokenID)
		pipe.Del(ctx, key)
		pipe.Del(ctx, buildRefreshKey(tokenHash))
	}
	pipe.Del(ctx, indexKey)

	_, err = pipe.Exec(ctx)
	return err
}

// GetRefreshTokenId Через хэшированный refresh token получаю tokenId чтобы потом искать по ИД ключу в списке
func (c *sessionCache) GetRefreshTokenId(ctx context.Context, hashedRefreshToken string) (string, error) {
	data, err := c.redis.Get(ctx, buildRefreshKey(hashedRefreshToken)).Result()
	return data, err
}

func (c *sessionCache) SetRefreshTokenId(ctx context.Context, hashedRefreshToken string, refreshTokenID string, ttl time.Duration) error {
	return c.redis.Set(ctx, buildRefreshKey(hashedRefreshToken), refreshTokenID, ttl).Err()
}

func (c *sessionCache) DeleteRefreshTokenId(ctx context.Context, hashedRefreshToken string) error {
	return c.redis.Del(ctx, buildRefreshKey(hashedRefreshToken)).Err()
}

// Формат ключа: auth:refresh:<userID>:<tokenID>
func buildSessionKey(tokenID string) string {
	return fmt.Sprintf("auth:refresh:%s", tokenID)
}

// Формат index ключа: auth:refresh:index:<userID>
func buildIndexKey(userID int64) string {
	return fmt.Sprintf("auth:refresh:index:%d", userID)
}

func buildRefreshKey(hashRefreshToken string) string {
	return fmt.Sprintf("auth:refresh_hash:%s", hashRefreshToken)
}
