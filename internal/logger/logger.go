package logger

import (
	"full-project-mock/internal/config"
	"log/slog"
	"os"
)

func InitLogger(config *config.Logger) (*slog.Logger, error) {
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

	return slog.New(loggerHandler), nil
}
