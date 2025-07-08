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

	cacheSession.On("GetRefreshTokenId", ctx, hashed).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteAllUserSessions", ctx, refreshSession.UserID).Return(nil)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.LogoutAllDevices(ctx, plainToken, clientIP, clientUserAgent)
	assert.NoError(t, err)

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashed)
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	cacheSession.AssertCalled(t, "DeleteAllUserSessions", ctx, refreshSession.UserID)
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

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
}

func TestLogoutAllDevices_GetSessionError(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.LogoutAllDevices(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
}

func TestLogoutAllDevices_DeleteAllUserSessions(t *testing.T) {
	ctx := context.Background()
	cacheSession := new(MockSessionCache)
	refreshSession := initRefreshSession()

	cacheSession.On("GetRefreshTokenId", ctx, hashRefreshToken(plainToken)).Return(refreshtokenId, nil)
	cacheSession.On("GetSession", ctx, refreshtokenId).Return(refreshSession, nil)
	cacheSession.On("DeleteAllUserSessions", ctx, refreshSession.UserID).Return(customErr)

	uc := Usecase{
		SessionCache: cacheSession,
	}

	err := uc.LogoutAllDevices(ctx, plainToken, clientIP, clientUserAgent)
	assert.EqualError(t, err, customErr.Error())

	cacheSession.AssertCalled(t, "GetRefreshTokenId", ctx, hashRefreshToken(plainToken))
	cacheSession.AssertCalled(t, "GetSession", ctx, refreshtokenId)
	cacheSession.AssertCalled(t, "DeleteAllUserSessions", ctx, refreshSession.UserID)
}
