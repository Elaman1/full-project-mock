package database

import (
	"fmt"
	"full-project-mock/internal/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
