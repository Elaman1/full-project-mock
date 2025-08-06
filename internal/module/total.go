package module

import (
	"database/sql"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"github.com/Elaman1/full-project-mock/internal/metrics"
	"github.com/Elaman1/full-project-mock/internal/module/user"
	"github.com/redis/go-redis/v9"
)

type Modules struct {
	UserHandler *user.UserHandler
}

// InitAllModule Инициализируем все модули здесь, Если новые добавиться то просто здесь же добавляем
func InitAllModule(db *sql.DB, redisDB *redis.Client, tokenService usecase.TokenService, metricsCollector metrics.MetricsCollector) *Modules {
	userHandler := user.InitUserModule(
		user.UserApp{DB: db, TokenService: tokenService, RedisDB: redisDB, MetricCollector: metricsCollector},
	)

	return &Modules{
		UserHandler: userHandler,
	}
}
