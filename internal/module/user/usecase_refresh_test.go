package user

import (
	"context"
	"full-project-mock/internal/domain/cache"
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
	tokenService := new(MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", user).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshtokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, hashRefreshToken(plainToken), refreshtokenId, mock.Anything).Return(nil)
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

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	tokenService.AssertCalled(t, "GenerateAccessToken", user)
	tokenService.AssertCalled(t, "GenerateRefreshToken")
	cacheSession.AssertCalled(t, "SetRefreshTokenId", ctx, hashRefreshToken(plainToken), refreshtokenId, mock.Anything)
	cacheSession.AssertCalled(t, "SaveSession", ctx, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything)
}

func TestRefresh_ParseTokenError(t *testing.T) {
	ctx := context.Background()
	tokenService := new(MockTokenService)
	regClaims := initRegisteredClaims()

	tokenService.On("ParseToken", accessToken).Return(regClaims, customErr)

	uc := Usecase{
		TokenService: tokenService,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertCalled(t, "ParseToken", accessToken)
}

func TestRefresh_GetById(t *testing.T) {
	ctx := context.Background()
	tokenService := new(MockTokenService)
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

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
}

func TestRefresh_GetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	tokenService := new(MockTokenService)
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
	cacheSession.On("GetRefreshTokenId", ctx, hashed).Return(refreshtokenId, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashed)
}

func TestRefresh_GetSessionError(t *testing.T) {
	ctx := context.Background()
	tokenService := new(MockTokenService)
	mockRepo := new(MockUserRepository)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(nil, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
}

func TestRefresh_GenerateAccessTokenError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
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

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
}

func TestRefresh_GenerateRefreshTokenError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshtokenId, plainToken, customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
	tokenService.AssertCalled(t, "GenerateRefreshToken")
}

func TestRefresh_SetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)
	cacheSession := new(MockSessionCache)

	regClaims := initRegisteredClaims()
	refreshSession := initRefreshSession()
	user, err := initUserWithPassword()
	if err != nil {
		t.Fatal(err)
	}

	tokenService.On("ParseToken", accessToken).Return(regClaims, nil)
	mockRepo.On("GetById", ctx, int64(defaultUserId)).Return(user, nil)
	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshtokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, hashRefreshToken(plainToken), refreshtokenId, mock.Anything).Return(customErr)

	uc := Usecase{
		TokenService: tokenService,
		Rep:          mockRepo,
		SessionCache: cacheSession,
	}

	accToken, plnToken, err := uc.Refresh(ctx, accessToken, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())
	assert.Equal(t, accToken, "")
	assert.Equal(t, plnToken, "")

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
	tokenService.AssertCalled(t, "GenerateRefreshToken")
	cacheSession.AssertCalled(t, "SetRefreshTokenId", ctx, hashRefreshToken(plainToken), refreshtokenId, mock.Anything)
}

func TestRefresh_SaveSessionError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	tokenService := new(MockTokenService)
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
	cacheSession.On("GetRefreshTokenId", ctx, hashed).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	tokenService.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshtokenId, plainToken, nil)
	cacheSession.On("SetRefreshTokenId", ctx, hashed, refreshtokenId, mock.Anything).Return(nil)
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

	tokenService.AssertCalled(t, "ParseToken", accessToken)
	mockRepo.AssertCalled(t, "GetById", ctx, int64(defaultUserId))
	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashed)
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	tokenService.AssertCalled(t, "GenerateAccessToken", mock.AnythingOfType("*model.User"))
	tokenService.AssertCalled(t, "GenerateRefreshToken")
	cacheSession.AssertCalled(t, "SetRefreshTokenId", ctx, hashed, refreshtokenId, mock.Anything)
	cacheSession.AssertCalled(t, "SaveSession", ctx, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything)
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
		TokenID:   refreshtokenId,
		TokenHash: hashRefreshToken(plainToken),
		UserAgent: clientUserAgent,
		IP:        clientIP,
	}
}
