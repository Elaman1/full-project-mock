package app

import (
	"context"
	"errors"
	"full-project-mock/internal/bootstrap"
	"full-project-mock/internal/config"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func RunApp() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig("./config/config.yaml", ".env")
	if err != nil {
		return err
	}

	app, err := bootstrap.InitApp(ctx, cfg)
	if err != nil {
		return err
	}

	errChan := make(chan error)

	go func() {
		app.Logger.Info("Starting server...", slog.String("addr", app.Server.Addr))
		if srvErr := app.Server.ListenAndServe(); srvErr != nil && !errors.Is(http.ErrServerClosed, srvErr) {
			app.Logger.Error("Server error", slog.Any("error", srvErr))
			errChan <- srvErr
		}
	}()

	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		app.Logger.Info("Shutdown initiated...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = app.Server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error("Graceful shutdown failed", slog.Any("error", err))
		return err
	}

	// Закрываем БД
	if app.DB != nil {
		err = app.DB.Close()
		if err != nil {
			app.Logger.Error("DB close failed", slog.Any("error", err))
			return err
		}
	}

	if app.RedisDB != nil {
		err = app.RedisDB.Close()
		if err != nil {
			app.Logger.Error("RedisDB close failed", slog.Any("error", err))
			return err
		}
	}

	app.Logger.Info("Shutdown complete.")
	return nil
}
