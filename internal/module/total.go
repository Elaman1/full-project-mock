package module

import (
	"database/sql"
	"full-project-mock/internal/domain/usecase"
	"full-project-mock/internal/module/user"
	"github.com/redis/go-redis/v9"
)

type Modules struct {
	UserHandler *user.UserHandler
}

// InitAllModule Инициализируем все модули здесь, Если новые добавиться то просто здесь же добавляем
func InitAllModule(db *sql.DB, redisDB *redis.Client, tokenService usecase.TokenService) *Modules {
	userHandler := user.InitUserModule(db, redisDB, tokenService)
	return &Modules{
		UserHandler: userHandler,
	}
}
