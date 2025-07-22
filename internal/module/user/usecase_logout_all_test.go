package user

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogoutAllDevices_Success(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()
	hashed := hashRefreshToken(plainToken)

	cacheSession.On("GetRefreshTokenId", ctx, hashed).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteAllUserSessions", ctx, refreshSession.UserID).Return(nil)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.LogoutAllDevices(ctx, plainToken, clientIP, clientUserAgent)
	assert.NoError(t, err)

	cacheSession.AssertExpectations(t)
}

func TestLogoutAllDevices_GetRefreshTokenIdError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return("", customErr)
	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.LogoutAllDevices(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertExpectations(t)
}

func TestLogoutAllDevices_GetSessionError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.LogoutAllDevices(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertExpectations(t)
}

func TestLogoutAllDevices_DeleteAllUserSessions(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshTokenId, nil)
	cacheSession.On("GetSession", ctx, refreshTokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteAllUserSessions", ctx, refreshSession.UserID).Return(customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.LogoutAllDevices(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertExpectations(t)
}
