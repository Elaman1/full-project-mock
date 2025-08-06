package middleware

import (
	"github.com/Elaman1/full-project-mock/internal/metrics"
	"log"
	"net/http"
	"time"
)

func MetricsMiddleware(mc metrics.MetricsCollector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//// Пропускаем вход на страницу метрики, их не надо считывать
			//if r.URL.Path == "/metrics" {
			//	next.ServeHTTP(w, r)
			//	return
			//}

			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()
			mc.IncHttpRequest(r.Method, r.URL.Path, http.StatusText(rw.statusCode))
			log.Println("Duration observed:", duration)
			mc.ObserveRequestDuration(r.Method, r.URL.Path, duration)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
