package user

import (
	"context"
	"fmt"
	domcache "full-project-mock/internal/domain/cache"
	"full-project-mock/internal/domain/constants"
	"full-project-mock/internal/domain/model"
	"full-project-mock/internal/domain/repository"
	"full-project-mock/internal/domain/usecase"
	"full-project-mock/pkg/hasher"
	"time"
)

type Usecase struct {
	Rep          repository.UserRepository
	TokenService usecase.TokenService
	SessionCache domcache.SessionCache
	RefreshTtl   time.Duration
}

func NewUserUsecase(userRepository repository.UserRepository, tokenService usecase.TokenService, sessionCache domcache.SessionCache) usecase.UserUsecase {
	return &Usecase{
		Rep:          userRepository,
		TokenService: tokenService,
		SessionCache: sessionCache,
		RefreshTtl:   7 * 24 * time.Hour, // 7 дней
	}
}

func (u *Usecase) Register(ctx context.Context, email, username, password string) (int64, error) {
	exists, err := u.Rep.Exists(ctx, email)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("пользователь с таким email %s уже существует", email)
	}

	pwd, err := hasher.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("произошла ошибка при хешировании пароля")
	}

	user := &model.User{
		Username: username,
		Password: pwd,
		Email:    email,
		RoleID:   constants.DefaultUserRoleID,
	}

	err = u.Rep.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("произошла ошибка при регистрации")
	}

	// пока что поставим 0, если потом надо будет вернем нормально
	return 0, nil
}

func (u *Usecase) Login(ctx context.Context, email, password, clientIP, ua string) (string, string, error) {
	user, err := u.Rep.Get(ctx, email)
	if err != nil {
		return "", "", err
	}

	err = hasher.Verify(user.Password, password)
	if err != nil {
		return "", "", err
	}

	accessToken, err := u.TokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, plain, err := u.TokenService.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	sess := &domcache.RefreshSession{
		UserID:    user.ID,
		TokenID:   refreshToken,
		TokenHash: hasher.Sha256Hex(plain),
		ExpiresAt: time.Now().Add(u.RefreshTtl),
		IP:        clientIP,
		UserAgent: ua,
	}

	err = u.SessionCache.SaveSession(ctx, sess, u.RefreshTtl)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
