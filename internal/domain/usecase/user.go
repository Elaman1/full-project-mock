package usecase

import "context"

type UserUsecase interface {
	Register(ctx context.Context, email, username, password string) (int64, error)
	Login(ctx context.Context, email, password, clientIP, ua string) (string, string, int, error)
	Refresh(ctx context.Context, accessToken, refreshToken, clientIP, ua string) (string, string, int, error)
	Logout(ctx context.Context, refreshToken, clientIP, ua string) error
	LogoutAllDevices(ctx context.Context, refreshToken, clientIP, ua string) error
}
