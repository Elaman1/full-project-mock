package logger

import (
	"full-project-mock/internal/config"
	"log/slog"
	"os"
	"sync"
)

var (
	logInstance *slog.Logger
	once        sync.Once
)

func InitLogger(config *config.Logger) *slog.Logger {
	// На всякий случаи, но лучше не вызывать этот метод кроме bootstrap
	once.Do(func() {
		var loggerHandler slog.Handler

		switch config.Format {
		case "json":
			loggerHandler = slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.Level(config.Level),
				})
		case "text":
			loggerHandler = slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.Level(config.Level),
				})
		}
		logInstance = slog.New(loggerHandler)
	})

	return logInstance
}
