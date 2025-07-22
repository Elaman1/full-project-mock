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

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteSession", ctx, refreshSession.UserID, refreshTokenId).Return(nil)
	cacheSession.On("DeleteRefreshTokenId", ctx, refreshSession.TokenHash).Return(nil)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.NoError(t, err)

	cacheSession.AssertExpectations(t)
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

	cacheSession.AssertExpectations(t)
}

func TestLogout_GetSessionError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(nil, customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertExpectations(t)
}

func TestLogout_SetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteSession", ctx, refreshSession.UserID, refreshTokenId).Return(customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertExpectations(t)
}

func TestLogout_DeleteRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteSession", ctx, refreshSession.UserID, refreshTokenId).Return(nil)
	cacheSession.On("DeleteRefreshTokenId", ctx, refreshSession.TokenHash).Return(customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.Logout(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertExpectations(t)
}
