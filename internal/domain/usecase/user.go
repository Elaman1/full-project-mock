package usecase

import "context"

type UserUsecase interface {
	Register(ctx context.Context, email, username, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, string, error)
}
