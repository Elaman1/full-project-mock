package cache

import (
	"context"
	"time"
)

type SessionCache interface {
	StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error
}
