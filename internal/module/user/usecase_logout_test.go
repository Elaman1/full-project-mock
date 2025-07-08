package user

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogout_Success(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteSession", ctx, refreshSession.UserID, refreshtokenId).Return(nil)
	cacheSession.On("DeleteRefreshTokenId", ctx, refreshSession.TokenHash).Return(nil)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.NoError(t, err)

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	cacheSession.AssertCalled(t, "DeleteSession", ctx, refreshSession.UserID, refreshtokenId)
	cacheSession.AssertCalled(t, "DeleteRefreshTokenId", ctx, refreshSession.TokenHash)
}

func TestLogout_GetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return("", customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
}

func TestLogout_GetSessionError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(nil, customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
}

func TestLogout_SetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteSession", ctx, refreshSession.UserID, refreshtokenId).Return(customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	cacheSession.AssertCalled(t, "DeleteSession", ctx, refreshSession.UserID, refreshtokenId)
}

func TestLogout_DeleteRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteSession", ctx, refreshSession.UserID, refreshtokenId).Return(nil)
	cacheSession.On("DeleteRefreshTokenId", ctx, refreshSession.TokenHash).Return(customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	cacheSession.AssertCalled(t, "DeleteSession", ctx, refreshSession.UserID, refreshtokenId)
	cacheSession.AssertCalled(t, "DeleteRefreshTokenId", ctx, refreshSession.TokenHash)
}
