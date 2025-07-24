package user

import (
	"database/sql"
	"github.com/Elaman1/full-project-mock/internal/cache"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"github.com/redis/go-redis/v9"
)

func InitUserModule(db *sql.DB, redisDB *redis.Client, tokenService usecase.TokenService) *UserHandler {
	sessionCache := cache.NewSessionRedisRepository(redisDB)
	userRepo := NewUserRepository(db)
	userUsecase := NewUserUsecase(userRepo, tokenService, sessionCache)
	return NewUserHandler(userUsecase)
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: usecase}
}
