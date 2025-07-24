package middleware

import (
	"context"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const UserIDKey = contextKey("userID")

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

func SetUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
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

			if mapClaims.Subject == "" {
				http.Error(w, "invalid token: empty subject", http.StatusUnauthorized)
				return
			}

			if mapClaims.ExpiresAt != nil && mapClaims.ExpiresAt.Time.Before(time.Now()) {
				http.Error(w, "expired token", http.StatusUnauthorized)
				return
			}

			ctx := SetUserIDToContext(r.Context(), mapClaims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
