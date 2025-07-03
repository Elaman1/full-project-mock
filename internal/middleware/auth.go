package middleware

import (
	"context"
	"full-project-mock/internal/domain/usecase"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userIDKey = contextKey("userID")

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func AuthMiddleware(tokenSvc usecase.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			mapClaims, err := tokenSvc.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			if mapClaims.ExpiresAt != nil && mapClaims.ExpiresAt.Time.Before(time.Now()) {
				http.Error(w, "expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, mapClaims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
