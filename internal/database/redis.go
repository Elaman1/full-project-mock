package database

import (
	"context"
	"fmt"
	"full-project-mock/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitRedis(ctx context.Context, cfg *config.Redis) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	status := r.Ping(ctx)
	if err := status.Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return r, nil
}
