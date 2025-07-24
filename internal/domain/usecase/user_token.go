package usecase

import (
	"github.com/Elaman1/full-project-mock/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(user *model.User) (string, error)
	GenerateRefreshToken() (tokenID, plainToken string, err error)
	ParseToken(tokenStr string) (jwt.RegisteredClaims, error)
}
