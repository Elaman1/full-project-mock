package cache

import (
	"context"
	"time"
)

type SessionCache interface {
	SaveSession(ctx context.Context, s *RefreshSession, ttl time.Duration) error
	GetSession(ctx context.Context, tokenID string) (*RefreshSession, error)
	DeleteSession(ctx context.Context, userID int64, tokenID string) error
	DeleteAllUserSessions(ctx context.Context, userID int64) error
	GetRefreshTokenId(ctx context.Context, hashedRefreshToken string) (string, error)
	SetRefreshTokenId(ctx context.Context, hashedRefreshToken string, refreshTokenID string, ttl time.Duration) error
	DeleteRefreshTokenId(ctx context.Context, hashedRefreshToken string) error
}

type RefreshSession struct {
	UserID    int64     `json:"user_id"`              // кому принадлежит токен
	TokenID   string    `json:"token_id"`             // уникальный ID (UUID)
	TokenHash string    `json:"token_hash"`           // хеш самого refresh токена (для безопасности)
	ExpiresAt time.Time `json:"expires_at"`           // когда истечёт
	IP        string    `json:"ip,omitempty"`         // (опционально, по безопасности)
	UserAgent string    `json:"user_agent,omitempty"` // (опционально, по безопасности)
}
