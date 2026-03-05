package user

import (
	"context"
	"database/sql"
	"github.com/Elaman1/full-project-mock/internal/auditlogger"
	"github.com/Elaman1/full-project-mock/internal/cache"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"github.com/Elaman1/full-project-mock/internal/metrics"
	"github.com/redis/go-redis/v9"
)

func InitUserModule(appCtx context.Context, userApp UserApp) *UserHandler {
	sessionCache := cache.NewSessionRedisRepository(userApp.RedisDB)
	userRepo := NewUserRepository(userApp.DB)

	// Не готово еще
	authLogger := auditlogger.InitAuditLogger(appCtx)

	userUsecase := NewUserUsecase(userRepo, userApp.TokenService, sessionCache, authLogger)
	// Передаем сверху метрику так как не обязательно внутр засовывать
	return NewUserHandler(userUsecase, userApp.MetricCollector)
}

type UserApp struct {
	DB              *sql.DB
	RedisDB         *redis.Client
	TokenService    usecase.TokenService
	MetricCollector metrics.MetricsCollector
}

func NewUserHandler(usecase usecase.UserUsecase, metricsCollector metrics.MetricsCollector) *UserHandler {
	return &UserHandler{
		Usecase:         usecase,
		MetricCollector: metricsCollector,
	}
}
