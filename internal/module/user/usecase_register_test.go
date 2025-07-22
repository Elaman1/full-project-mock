package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegister(t *testing.T) {
	type testCase struct {
		name       string
		setupMocks func(*MockUserRepository)
		wantID     int64
		wantErr    error
		wantErrMsg string
	}

	cases := []testCase{
		{
			name: "success",
			setupMocks: func(repo *MockUserRepository) {
				repo.On("Exists", mock.Anything, defaultEmail).Return(false, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantID:  0,
			wantErr: nil,
		},
		{
			name: "email already exists",
			setupMocks: func(repo *MockUserRepository) {
				repo.On("Exists", mock.Anything, defaultEmail).Return(true, nil)
			},
			wantID:     0,
			wantErrMsg: fmt.Sprintf("пользователь с таким email %s уже существует", defaultEmail),
		},
		{
			name: "exists returns error",
			setupMocks: func(repo *MockUserRepository) {
				repo.On("Exists", mock.Anything, defaultEmail).Return(false, customErr)
			},
			wantID:     0,
			wantErrMsg: customErr.Error(),
		},
		{
			name: "create returns error",
			setupMocks: func(repo *MockUserRepository) {
				repo.On("Exists", mock.Anything, defaultEmail).Return(false, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("customer error"))
			},
			wantID:     0,
			wantErrMsg: "произошла ошибка при регистрации",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := new(MockUserRepository)
			tc.setupMocks(repo)

			uc := &Usecase{
				Rep: repo,
			}

			gotID, err := uc.Register(context.Background(), defaultEmail, defaultUserName, defaultPassword)

			assert.Equal(t, tc.wantID, gotID)

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else if tc.wantErrMsg != "" {
				assert.EqualError(t, err, tc.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
