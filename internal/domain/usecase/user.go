package usecase

import "context"

type UserUsecase interface {
	Register(ctx context.Context, email, username, password string) (int64, error)
	Login(ctx context.Context, email, password, clientIP, ua string) (string, string, error)
	Refresh(ctx context.Context, accessToken, refreshToken, clientIP, ua string) (string, string, error)
}
