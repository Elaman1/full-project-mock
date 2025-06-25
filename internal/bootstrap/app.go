package bootstrap

import (
	"context"
	"database/sql"
	"full-project-mock/internal/config"
	"full-project-mock/internal/database"
	"full-project-mock/internal/delivery/rest"
	"full-project-mock/internal/logger"
	"full-project-mock/internal/module"
	"full-project-mock/internal/service"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	Server *http.Server
	DB     *sql.DB
	Logger *slog.Logger
}

func InitApp(ctx context.Context, cfg *config.Config) (*App, error) {
	logs, err := logger.InitLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	db, err := database.InitPostgres(&cfg.PostgresDB)
	if err != nil {
		logs.Error("error connecting to the database", "error", err)
		return nil, err
	}

	redisDB := database.InitRedis(&cfg.Redis)

	tokenService := service.NewTokenService(cfg.JWT.Secret, time.Minute*15)
	allModules := module.InitAllModule(db, redisDB, tokenService)

	routeApp := &rest.RouteApp{
		Logs:         logs,
		TokenService: tokenService,
	}
	routeHandler := rest.InitRouter(ctx, routeApp, allModules)

	srv := &http.Server{
		Addr:         cfg.Server.Port,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		Handler:      routeHandler,
	}

	return &App{
		Server: srv,
		DB:     db, // Передаем, чтобы закрыть соединение при отключении сервера
		Logger: logs,
	}, nil
}
