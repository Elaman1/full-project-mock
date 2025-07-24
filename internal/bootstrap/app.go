package bootstrap

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"github.com/Elaman1/full-project-mock/internal/config"
	"github.com/Elaman1/full-project-mock/internal/database"
	"github.com/Elaman1/full-project-mock/internal/delivery/rest"
	"github.com/Elaman1/full-project-mock/internal/logger"
	"github.com/Elaman1/full-project-mock/internal/module"
	"github.com/Elaman1/full-project-mock/internal/service"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type App struct {
	Server  *http.Server
	DB      *sql.DB
	Logger  *slog.Logger
	RedisDB *redis.Client
}

func InitApp(ctx context.Context, cfg *config.Config) (*App, error) {
	logs := logger.InitLogger(&cfg.Logger)

	db, err := database.InitPostgres(&cfg.PostgresDB)
	if err != nil {
		logs.Error("error connecting to the database", "error", err)
		return nil, err
	}

	redisDB, err := database.InitRedis(ctx, &cfg.Redis)
	if err != nil {
		logs.Error("error connecting to the redis", "error", err)
		return nil, err
	}

	publicKey, err := LoadRSAPublicKey(cfg.JWT.PublicKeyPath)
	if err != nil {
		logs.Error("error loading public key", "error", err)
		return nil, err
	}

	privateKey, err := LoadRSAPrivateKey(cfg.JWT.PrivateKeyPath)
	if err != nil {
		logs.Error("error loading private key", "error", err)
		return nil, err
	}

	ttl, err := time.ParseDuration(cfg.JWT.AccessTTL)
	// По идее дополнительно сверху проверяется
	if err != nil {
		logs.Error("error parsing access ttl", "error", err)
		return nil, err
	}

	tokenService := service.NewTokenService(publicKey, privateKey, ttl)
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
		Server:  srv,
		DB:      db, // Передаем, чтобы закрыть соединение при отключении сервера
		Logger:  logs,
		RedisDB: redisDB, // То же самое
	}, nil
}

func LoadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid PEM format for private key")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("parsed key is not an RSA private key")
	}

	return rsaKey, nil
}

func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid PEM format for public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not an RSA public key")
	}

	return rsaPub, nil
}
