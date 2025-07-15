package config

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

func validateCfg(cfg *Config) error {
	if err := validateServer(cfg); err != nil {
		return err
	}

	if err := validatePostgresDB(cfg); err != nil {
		return err
	}

	if err := validateLogger(cfg); err != nil {
		return err
	}

	if err := validateRedis(cfg); err != nil {
		return err
	}

	if err := validateJWT(cfg); err != nil {
		return err
	}

	return nil
}

func validateJWT(cfg *Config) error {
	if err := validateJWTAccessTTL(cfg); err != nil {
		return err
	}

	if cfg.JWT.PrivateKeyPath == "" {
		return errors.New("private key file path is required")
	}

	if cfg.JWT.PublicKeyPath == "" {
		return errors.New("public key file path is required")
	}

	return nil
}

func validateLogger(cfg *Config) error {
	if err := validateLogLevel(cfg.Logger.Level); err != nil {
		return err
	}

	if cfg.Logger.Format == "" {
		return errors.New("logger no format specified")
	}

	return nil
}

func validateLogLevel(level int) error {
	switch level {
	case int(slog.LevelDebug),
		int(slog.LevelInfo),
		int(slog.LevelWarn),
		int(slog.LevelError):
		return nil
	default:
		return fmt.Errorf("invalid log level: %d", level)
	}
}

func validateServer(cfg *Config) error {
	if cfg.Server.Port == "" {
		return errors.New("missing required configuration: server_port")
	}

	if cfg.Server.ReadTimeout == 0 {
		return errors.New("missing required configuration variable: server_read_timeout")
	}

	if cfg.Server.WriteTimeout == 0 {
		return errors.New("missing required configuration variable: server_write_timeout")
	}

	return nil
}

func validatePostgresDB(cfg *Config) error {
	if cfg.PostgresDB.Host == "" {
		return errors.New("missing required configuration variable: postgres_host")
	}

	if cfg.PostgresDB.Port == 0 {
		return errors.New("missing required configuration variable: postgres_port")
	}

	if cfg.PostgresDB.DBName == "" {
		return errors.New("missing required configuration variable: postgres_dbname")
	}

	if cfg.PostgresDB.SslMode == "" {
		return errors.New("missing required configuration variable: postgres_ssl_mode")
	}

	if cfg.PostgresDB.MaxOpenConns == 0 {
		return errors.New("missing required configuration variable: postgres_max_open_conns")
	}

	if cfg.PostgresDB.MaxIdleConns == 0 {
		return errors.New("missing required configuration variable: postgres_max_idle_conns")
	}

	if cfg.PostgresDB.MaxLifeTime == 0 {
		return errors.New("missing required configuration variable: postgres_max_life_time")
	}

	return nil
}

func validateRedis(cfg *Config) error {
	if cfg.Redis.Host == "" {
		return errors.New("missing required configuration variable: redis_host")
	}

	if cfg.Redis.Port == 0 {
		return errors.New("missing required configuration variable: redis_port")
	}

	return nil
}

func validateJWTAccessTTL(cfg *Config) error {
	if cfg.JWT.AccessTTL == "" {
		return errors.New("missing required configuration variable: jwt_access_ttl")
	}

	_, err := time.ParseDuration(cfg.JWT.AccessTTL)
	if err != nil {
		return err
	}

	return nil
}
