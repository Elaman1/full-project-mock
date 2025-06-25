package usecase

import "full-project-mock/internal/domain/model"

type TokenService interface {
	GenerateAccessToken(user *model.User) (string, error)
	GenerateRefreshToken() (string, error)
	ParseToken(tokenStr string) (string, error) // вернёт userID (sub)
}
