package user

import (
	"context"
	"errors"
	"fmt"
	"full-project-mock/internal/domain/cache"
	"full-project-mock/internal/domain/model"
	"full-project-mock/pkg/hasher"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"strconv"
	"time"
)

var (
	defaultEmail    = "test@example.com"
	defaultUserName = "testuser"
	defaultPassword = "securepassword"
	clientUserAgent = "testClient"
	clientIP        = "127.0.0.1"

	accessToken    = "testAccessToken"
	refreshTokenId = "testTokenId"
	plainToken     = "testPlainToken"

	defaultUserId = 11

	customErr = errors.New("custom error")
)

// Session раздел
type MockSessionCache struct {
	mock.Mock
}

func (m *MockSessionCache) SaveSession(ctx context.Context, s *cache.RefreshSession, ttl time.Duration) error {
	args := m.Called(ctx, s, ttl)
	return args.Error(0)
}

func (m *MockSessionCache) GetSession(ctx context.Context, tokenID string) (*cache.RefreshSession, error) {
	args := m.Called(ctx, tokenID)
	c := args.Get(0)
	if c == nil {
		return nil, args.Error(1)
	}

	refreshSession, ok := c.(*cache.RefreshSession)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return refreshSession, args.Error(1)
}

func (m *MockSessionCache) DeleteSession(ctx context.Context, userID int64, tokenID string) error {
	args := m.Called(ctx, userID, tokenID)
	return args.Error(0)
}

func (m *MockSessionCache) DeleteAllUserSessions(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionCache) GetRefreshTokenId(ctx context.Context, hashedRefreshToken string) (string, error) {
	args := m.Called(ctx, hashedRefreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockSessionCache) SetRefreshTokenId(ctx context.Context, hashedRefreshToken string, refreshTokenID string, ttl time.Duration) error {
	args := m.Called(ctx, hashedRefreshToken, refreshTokenID, ttl)
	return args.Error(0)
}

func (m *MockSessionCache) DeleteRefreshTokenId(ctx context.Context, hashedRefreshToken string) error {
	args := m.Called(ctx, hashedRefreshToken)
	return args.Error(0)
}

func initUserWithPassword() (*model.User, error) {
	password, err := hasher.HashPassword(defaultPassword)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Password: password,
		ID:       int64(defaultUserId),
	}

	return &user, nil
}

// Refresh раздел
func initRegisteredClaims() jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Subject: strconv.Itoa(defaultUserId),
	}
}

func initRefreshSession() *cache.RefreshSession {
	return &cache.RefreshSession{
		UserID:    int64(defaultUserId),
		ExpiresAt: time.Now().Add(time.Minute * 10),
		TokenID:   refreshTokenId,
		TokenHash: hashRefreshToken(plainToken),
		UserAgent: clientUserAgent,
		IP:        clientIP,
	}
}

// MockUserRepository раздел
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Get(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}

	user, ok := u.(*model.User)
	if !ok {
		return nil, fmt.Errorf("error casting model.User")
	}

	return user, args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	user, ok := args.Get(0).(*model.User)
	if !ok {
		return nil, fmt.Errorf("error casting model.User")
	}

	return user, args.Error(1)
}
