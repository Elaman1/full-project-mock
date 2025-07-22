package user

import (
	"context"
	"full-project-mock/internal/domain/cache"
	"full-project-mock/internal/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
	"time"
)

func TestRefresh_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(mocks.MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", user).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshTokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, hashRefreshToken(plainToken), refreshTokenId, mock.Anything).Return(nil)
	cacheSession.On("SaveSession", ctx, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything).Return(nil)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.NoError(t, err)
	assert.Equal(t, accToken, accessToken)
	assert.Equal(t, plnToken, plainToken)

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	cacheSession.AssertExpectations(t)
}

func TestRefresh_ParseTokenError(t *testing.T) {
	ctx := context.Background()
	tokenService := new(mocks.MockTokenService)
	regClaims := initRegisteredClaims()

	tokenService.On("ParseToken", accessToken).Return(regClaims, customErr)

	uc := Usecase{
		TokenService: tokenService,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
}

func TestRefresh_GetById(t *testing.T) {
	ctx := context.Background()
	tokenService := new(mocks.MockTokenService)
	mockRepo := new(MockUserRepository)

	regClaims := initRegisteredClaims()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestRefresh_GetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	tokenService := new(mocks.MockTokenService)
	mockRepo := new(MockUserRepository)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	hashed := hashRefreshToken(plainToken)
	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashed).Return(refreshTokenId, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	cacheSession.AssertExpectations(t)
}

func TestRefresh_GetSessionError(t *testing.T) {
	ctx := context.Background()
	tokenService := new(mocks.MockTokenService)
	mockRepo := new(MockUserRepository)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(nil, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	cacheSession.AssertExpectations(t)
}

func TestRefresh_GenerateAccessTokenError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(mocks.MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	cacheSession.AssertExpectations(t)
}

func TestRefresh_GenerateRefreshTokenError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(mocks.MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshTokenId, plainToken, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	cacheSession.AssertExpectations(t)
}

func TestRefresh_SetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(mocks.MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshTokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, hashRefreshToken(plainToken), refreshTokenId, mock.Anything).Return(customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	cacheSession.AssertExpectations(t)
}

func TestRefresh_SaveSessionError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(mocks.MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	hashed := hashRefreshToken(plainToken)
	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashed).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshTokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, hashed, refreshTokenId, mock.Anything).Return(nil)
	cacheSession.On("SaveSession", ctx, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything).Return(customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	cacheSession.AssertExpectations(t)
}

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
