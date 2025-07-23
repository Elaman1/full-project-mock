package user

import (
	"context"
	"full-project-mock/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogin(t *testing.T) {
	user, err := initUserWithPassword()
	require.NoError(t, err)

	type testCase struct {
		name       string
		setupMocks func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache)
		wantToken  string
		wantPlain  string
		wantErr    error
	}

	cases := []testCase{
		{
			name: "success",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				repo.On("Get", mock.Anything, defaultEmail).
					Return(user, nil)
				ts.On("GenerateAccessToken", mock.Anything).
					Return(accessToken, nil)
				ts.On("GenerateRefreshToken").
					Return(refreshTokenId, plainToken, nil)
				cs.On("SetRefreshTokenId", mock.Anything, mock.Anything, refreshTokenId, mock.Anything).
					Return(nil)
				cs.On("SaveSession", mock.Anything, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything).
					Return(nil)
			},
			wantToken: accessToken,
			wantPlain: plainToken,
			wantErr:   nil,
		},
		{
			name: "repo returns error",
			setupMocks: func(repo *MockUserRepository, _ *mocks.MockTokenService, _ *MockSessionCache) {
				repo.On("Get", mock.Anything, defaultEmail).
					Return(nil, customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "generate access token error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, _ *MockSessionCache) {
				repo.On("Get", mock.Anything, defaultEmail).
					Return(user, nil)
				ts.On("GenerateAccessToken", mock.Anything).
					Return("", customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "generate refresh token error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, _ *MockSessionCache) {
				repo.On("Get", mock.Anything, defaultEmail).
					Return(user, nil)
				ts.On("GenerateAccessToken", mock.Anything).
					Return(accessToken, nil)
				ts.On("GenerateRefreshToken").
					Return("", "", customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "set refresh token id error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				repo.On("Get", mock.Anything, defaultEmail).
					Return(user, nil)
				ts.On("GenerateAccessToken", mock.Anything).
					Return(accessToken, nil)
				ts.On("GenerateRefreshToken").
					Return(refreshTokenId, plainToken, nil)
				cs.On("SetRefreshTokenId", mock.Anything, mock.Anything, refreshTokenId, mock.Anything).
					Return(customErr)
			},
			wantToken: "",
			wantPlain: "",
			wantErr:   customErr,
		},
		{
			name: "save session error",
			setupMocks: func(repo *MockUserRepository, ts *mocks.MockTokenService, cs *MockSessionCache) {
				repo.On("Get", mock.Anything, defaultEmail).
					Return(user, nil)
				ts.On("GenerateAccessToken", mock.Anything).
					Return(accessToken, nil)
				ts.On("GenerateRefreshToken").
					Return(refreshTokenId, plainToken, nil)
				cs.On("SetRefreshTokenId", mock.Anything, mock.Anything, refreshTokenId, mock.Anything).
					Return(nil)
				cs.On("SaveSession", mock.Anything, mock.AnythingOfType("*cache.RefreshSession"), mock.Anything).
					Return(customErr)
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
			tokenSvc := new(mocks.MockTokenService)
			cache := new(MockSessionCache)

			tc.setupMocks(repo, tokenSvc, cache)

			uc := Usecase{
				Rep:          repo,
				TokenService: tokenSvc,
				SessionCache: cache,
			}

			token, plain, err := uc.Login(context.Background(), defaultEmail, defaultPassword, clientIP, clientUserAgent)

			assert.Equal(t, tc.wantToken, token)
			assert.Equal(t, tc.wantPlain, plain)

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			tokenSvc.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}
