package middleware

import (
	"github.com/Elaman1/full-project-mock/internal/service"
	"log/slog"
	"net/http"
	"time"
)

func LogMiddleware(logs *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceId, err := service.GenerateTraceID()
			if err != nil {
				logs.Error(err.Error())
			}

			ctx := service.WithLogger(r.Context(), logs.With("traceId", traceId))

			srw := &StatusResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			start := time.Now()
			next.ServeHTTP(srw, r.WithContext(ctx))
			duration := time.Since(start)

			// Фиксируем только успешные, так как ошибочные будут отдельно логироваться
			if srw.statusCode != http.StatusOK {
				return
			}

			service.LoggerFromContext(ctx).Info("request completed",
				"path", r.URL.Path,
				"method", r.Method,
				"duration", duration,
			)
		})
	}
}
