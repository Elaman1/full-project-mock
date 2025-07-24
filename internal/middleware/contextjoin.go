package middleware

import (
	"context"
	"github.com/Elaman1/full-project-mock/internal/service"
	"net/http"
)

func ContextJoinMiddleware(shutdownCtx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			joinedCtx := service.JoinContexts(r.Context(), shutdownCtx)
			next.ServeHTTP(w, r.WithContext(joinedCtx))
		})
	}
}
