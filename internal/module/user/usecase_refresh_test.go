package user

import (
	"context"
	"github.com/Elaman1/full-project-mock/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRefresh(t *testing.T) {
	hashed := hashRefreshToken(plainToken)
	refreshSession := initRefreshSession()
	regClaims := initRegisteredClaims()
	user, err := initUserWithPassword()
	require.NoError(t, err)

	type testCase struct {
		name       string
		setupMocks func(*MockUserRepository, *mocks.MockTokenService, *MockSessionCache)
		wantToken  string
		wantPlain  string
		wantErr    error
	}

	cases := []testCase{
		{
			name: "success",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, nil)
				cs.On("GetRefreshTokenId", mock.Anything, hashed).Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).Return(refreshSession, nil)
				ts.On("GenerateAccessToken", user).Return(accessToken, nil)
				ts.On("GenerateRefreshToken").Return(refreshTokenId, plainToken, nil)
				cs.On("SetRefreshTokenId", mock.Anything, hashed, refreshTokenId, mock.Anything).Return(nil)
				cs.On("SaveSession", mock.Anything, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything).Return(nil)
			},
			wantToken: accessToken,
			wantPlain: plainToken,
			wantErr:   nil,
		},
		{
			name: "parse token returns error",
			setupMocks: func(_ *MockUserRepository, ts *mocks.MockTokenService, _ *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "get user by id error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, _ *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "get refresh token id error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, nil)
				cs.On("GetRefreshTokenId", mock.Anything, hashed).Return("", customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "get session error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, nil)
				cs.On("GetRefreshTokenId", mock.Anything, hashed).Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).Return(nil, customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "generate access token error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, nil)
				cs.On("GetRefreshTokenId", mock.Anything, hashed).Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).Return(refreshSession, nil)
				ts.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return("", customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "generate refresh token error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, nil)
				cs.On("GetRefreshTokenId", mock.Anything, hashed).Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).Return(refreshSession, nil)
				ts.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
				ts.On("GenerateRefreshToken").Return("", "", customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "set refresh token id error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, nil)
				cs.On("GetRefreshTokenId", mock.Anything, hashed).Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).Return(refreshSession, nil)
				ts.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
				ts.On("GenerateRefreshToken").Return(refreshTokenId, plainToken, nil)
				cs.On("SetRefreshTokenId", mock.Anything, hashed, refreshTokenId, mock.Anything).Return(customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "save session error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				ts.On("ParseToken", accessToken).Return(regClaims, nil)
				repo.On("GetById", mock.Anything, int64(defaultUserId)).Return(user, nil)
				cs.On("GetRefreshTokenId", mock.Anything, hashed).Return(refreshTokenId, nil)
				cs.On("GetSession", mock.Anything, refreshTokenId).Return(refreshSession, nil)
				ts.On("GenerateAccessToken", mock.AnythingOfType("*model.User")).Return(accessToken, nil)
				ts.On("GenerateRefreshToken").Return(refreshTokenId, plainToken, nil)
				cs.On("SetRefreshTokenId", mock.Anything, hashed, refreshTokenId, mock.Anything).Return(nil)
				cs.On("SaveSession", mock.Anything, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything).Return(customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := new(MockUserRepository)
			ts := new(mocks.MockTokenService)
			cs := new(MockSessionCache)

			tc.setupMocks(repo, ts, cs)

			uc := Usecase{
				Rep:          repo,
				TokenService: ts,
				SessionCache: cs,
			}

			gotToken, gotPlain, err := uc.Refresh(context.Background(), accessToken, plainToken, clientIP, clientUserAgent)

			assert.Equal(t, tc.wantToken, gotToken)
			assert.Equal(t, tc.wantPlain, gotPlain)

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			ts.AssertExpectations(t)
			cs.AssertExpectations(t)
		})
	}
}
