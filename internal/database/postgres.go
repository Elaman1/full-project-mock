package database

import (
	"database/sql"
	"fmt"
	"github.com/Elaman1/full-project-mock/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func InitPostgres(cfg *config.PostgresDB) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SslMode,
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)

	return conn, nil
}
