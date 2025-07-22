package mocks

import (
	"errors"
	"full-project-mock/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(user *model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken() (tokenID, plainToken string, err error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockTokenService) ParseToken(tokenStr string) (jwt.RegisteredClaims, error) {
	args := m.Called(tokenStr)
	jwtClaims, ok := args.Get(0).(jwt.RegisteredClaims)
	if !ok {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}

	return jwtClaims, args.Error(1)
}
