package user

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLogout(t *testing.T) {
	hashed := hashRefreshToken(plainToken)
	refreshSession := initRefreshSession()

	type testCase struct {
		name       string
		setupMocks func(cs *MockSessionCache)
		wantErr    error
	}

	cases := []testCase{
		{
			name: "success",
			setupMocks: func(cs *MockSessionCache) {
				cs.On("GetRefreshTokenId", mock.Anything, hashed).
					Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).
					Return(refreshSession, nil)
				cs.On("DeleteSession", mock.Anything, refreshSession.UserID, refreshTokenId).
					Return(nil)
				cs.On("DeleteRefreshTokenId", mock.Anything, refreshSession.TokenHash).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "get refresh token id returns error",
			setupMocks: func(cs *MockSessionCache) {
				cs.On("GetRefreshTokenId", mock.Anything, hashed).
					Return("", customErr)
			},
			wantErr: customErr,
		},
		{
			name: "get session returns error",
			setupMocks: func(cs *MockSessionCache) {
				cs.On("GetRefreshTokenId", mock.Anything, hashed).
					Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).
					Return(nil, customErr)
			},
			wantErr: customErr,
		},
		{
			name: "delete session returns error",
			setupMocks: func(cs *MockSessionCache) {
				cs.On("GetRefreshTokenId", mock.Anything, hashed).
					Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).
					Return(refreshSession, nil)
				cs.On("DeleteSession", mock.Anything, refreshSession.UserID, refreshTokenId).
					Return(customErr)
			},
			wantErr: customErr,
		},
		{
			name: "delete refresh token id returns error",
			setupMocks: func(cs *MockSessionCache) {
				cs.On("GetRefreshTokenId", mock.Anything, hashed).
					Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).
					Return(refreshSession, nil)
				cs.On("DeleteSession", mock.Anything, refreshSession.UserID, refreshTokenId).
					Return(nil)
				cs.On("DeleteRefreshTokenId", mock.Anything, refreshSession.TokenHash).
					Return(customErr)
			},
			wantErr: customErr,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cacheSession := new(MockSessionCache)
			tc.setupMocks(cacheSession)

			uc := Usecase{
				SessionCache: cacheSession,
			}

			err := uc.Logout(context.Background(), plainToken, clientIP, clientUserAgent)

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			cacheSession.AssertExpectations(t)
		})
	}
}
