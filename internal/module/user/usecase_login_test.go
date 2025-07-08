package user

import (
	"context"
	"errors"
	"full-project-mock/internal/domain/cache"
	"full-project-mock/internal/domain/model"
	"full-project-mock/pkg/hasher"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

var (
	clientUserAgent = "testClient"
	clientIP        = "127.0.0.1"

	accessToken    = "testAccessToken"
	refreshtokenId = "testTokenId"
	plainToken     = "testPlainToken"

	defaultUserId = 11

	customErr = errors.New("custom error")
)

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)
	cacheSession := new(MockSessionCache)

	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	mockRepo.On("Get", ctx, defaultEmail).Return(user, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshtokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, mock.Anything, refreshtokenId, mock.Anything).Return(nil)
	cacheSession.On("SaveSession", ctx, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything).Return(nil)

	uc := Usecase{
		Rep:          mockRepo,
		SessionCache: cacheSession,
		TokenService: tokenService,
	}

	ac, pl, err := uc.Login(ctx, defaultEmail, defaultPassword, clientIP, clientUserAgent)
	assert.NoError(t, err)
	assert.Equal(t, ac, accessToken)
	assert.Equal(t, pl, plainToken)

	mockRepo.AssertCalled(t, "Get", ctx, defaultEmail)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
	tokenService.AssertCalled(t, "GenerateRefreshToken")
	cacheSession.AssertCalled(t, "SetRefreshTokenId", ctx, mock.Anything, refreshtokenId, mock.Anything)
	cacheSession.AssertCalled(t, "SaveSession", ctx, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything)
}

func TestLogin_GetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)

	mockRepo.On("Get", ctx, defaultEmail).Return(nil, customErr)

	uc := Usecase{
		Rep: mockRepo,
	}

	ac, pl, err := uc.Login(ctx, defaultEmail, defaultPassword, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	assert.Equal(t, ac, "")
	assert.Equal(t, pl, "")

	mockRepo.AssertCalled(t, "Get", ctx, defaultEmail)
}

func TestLogin_GenerateAccessTokenError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)

	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	mockRepo.On("Get", ctx, defaultEmail).Return(user, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return("", customErr)

	uc := Usecase{
		Rep:          mockRepo,
		TokenService: tokenService,
	}

	ac, pl, err := uc.Login(ctx, defaultEmail, defaultPassword, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, ac, "")
	assert.Equal(t, pl, "")
	mockRepo.AssertCalled(t, "Get", ctx, defaultEmail)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
}

func TestLogin_GenerateRefreshTokenError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)

	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	mockRepo.On("Get", ctx, defaultEmail).Return(user, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return("", "", customErr)

	uc := Usecase{
		Rep:          mockRepo,
		TokenService: tokenService,
	}

	ac, pl, err := uc.Login(ctx, defaultEmail, defaultPassword, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, ac, "")
	assert.Equal(t, pl, "")

	mockRepo.AssertCalled(t, "Get", ctx, defaultEmail)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
	tokenService.AssertCalled(t, "GenerateRefreshToken")
}

func TestLogin_SetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)
	cacheSession := new(MockSessionCache)

	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	customErr = errors.New("custom error")

	mockRepo.On("Get", ctx, defaultEmail).Return(user, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshtokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, mock.Anything, refreshtokenId, mock.Anything).Return(customErr)

	uc := Usecase{
		Rep:          mockRepo,
		SessionCache: cacheSession,
		TokenService: tokenService,
	}

	ac, pl, err := uc.Login(ctx, defaultEmail, defaultPassword, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, ac, "")
	assert.Equal(t, pl, "")
	mockRepo.AssertCalled(t, "Get", ctx, defaultEmail)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
	tokenService.AssertCalled(t, "GenerateRefreshToken")
	cacheSession.AssertCalled(t, "SetRefreshTokenId", ctx, mock.Anything, refreshtokenId, mock.Anything)
}

func TestLogin_SaveSessionError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)
	cacheSession := new(MockSessionCache)

	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	mockRepo.On("Get", ctx, defaultEmail).Return(user, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshtokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, mock.Anything, refreshtokenId, mock.Anything).Return(nil)
	cacheSession.On("SaveSession", ctx, mock.Anything, mock.Anything).Return(customErr)

	uc := Usecase{
		Rep:          mockRepo,
		SessionCache: cacheSession,
		TokenService: tokenService,
	}

	ac, pl, err := uc.Login(ctx, defaultEmail, defaultPassword, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, ac, "")
	assert.Equal(t, pl, "")
	mockRepo.AssertCalled(t, "Get", ctx, defaultEmail)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
	tokenService.AssertCalled(t, "GenerateRefreshToken")
	cacheSession.AssertCalled(t, "SetRefreshTokenId", ctx, mock.Anything, refreshtokenId, mock.Anything)
	cacheSession.AssertCalled(t, "SaveSession", ctx, mock.Anything, mock.Anything)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(user *model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken() (tokenID, plainToken string, err error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockTokenService) ParseToken(tokenStr string) (jwt.RegisteredClaims, error) {
	args := m.Called(tokenStr)
	jwtClaims, ok := args.Get(0).(jwt.RegisteredClaims)
	if !ok {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}

	return jwtClaims, args.Error(1)
}

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
