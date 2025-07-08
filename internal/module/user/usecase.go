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
	"strconv"
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

	// пока что поставим 0, если потом надо, будем возвращать нормально
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

	return u.generateAccessAndRefreshToken(ctx, clientIP, ua, user)
}

func (u *Usecase) generateAccessAndRefreshToken(ctx context.Context, clientIP, ua string, user *model.User) (string, string, error) {
	accessToken, err := u.TokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshTokenId, plainToken, err := u.TokenService.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	hashedPlainToken := hashRefreshToken(plainToken)
	newSess := &domcache.RefreshSession{
		UserID:    user.ID,
		TokenID:   refreshTokenId,
		TokenHash: hashedPlainToken,
		ExpiresAt: time.Now().Add(u.RefreshTtl),
		IP:        clientIP,
		UserAgent: ua,
	}

	err = u.SessionCache.SetRefreshTokenId(ctx, hashedPlainToken, refreshTokenId, u.RefreshTtl)
	if err != nil {
		return "", "", err
	}

	err = u.SessionCache.SaveSession(ctx, newSess, u.RefreshTtl)
	if err != nil {
		return "", "", err
	}

	return accessToken, plainToken, nil
}

func (u *Usecase) Refresh(ctx context.Context, accessToken, refreshToken, clientIP, ua string) (string, string, error) {
	mapClaims, err := u.TokenService.ParseToken(accessToken)
	if err != nil {
		return "", "", err
	}

	userId, err := strconv.Atoi(mapClaims.Subject)
	if err != nil {
		return "", "", err
	}

	user, err := u.Rep.GetById(ctx, int64(userId))
	if err != nil {
		return "", "", err
	}

	refreshTokenId, err := u.SessionCache.GetRefreshTokenId(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		return "", "", err
	}

	refreshSession, err := u.SessionCache.GetSession(ctx, refreshTokenId)
	if err != nil {
		return "", "", err
	}

	if refreshSession.ExpiresAt.Before(time.Now()) {
		return "", "", fmt.Errorf("время истек заново авторизуйтесь")
	}

	if refreshSession.IP != clientIP {
		return "", "", fmt.Errorf("авторизуйтесь еще раз")
	}

	if refreshSession.UserAgent != ua {
		return "", "", fmt.Errorf("авторизуйтесь еще раз")
	}

	if hashRefreshToken(refreshToken) != refreshSession.TokenHash {
		return "", "", fmt.Errorf("авторизуйтесь еще раз")
	}

	return u.generateAccessAndRefreshToken(ctx, clientIP, ua, user)
}

func (u *Usecase) Logout(ctx context.Context, refreshToken, clientIP, ua string) error {
	refreshTokenId, err := u.SessionCache.GetRefreshTokenId(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		return err
	}

	refreshSession, err := u.SessionCache.GetSession(ctx, refreshTokenId)
	if err != nil {
		return err
	}

	if refreshSession.TokenHash != hashRefreshToken(refreshToken) {
		return fmt.Errorf("невозможно выполнить операцию")
	}

	if refreshSession.IP != clientIP {
		return fmt.Errorf("невозможно выполнить операцию")
	}

	if refreshSession.UserAgent != ua {
		return fmt.Errorf("невозможно выполнить операцию")
	}

	err = u.SessionCache.DeleteSession(ctx, refreshSession.UserID, refreshTokenId)
	if err != nil {
		return err
	}

	err = u.SessionCache.DeleteRefreshTokenId(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) LogoutAllDevices(ctx context.Context, refreshToken, clientIP, ua string) error {
	refreshTokenId, err := u.SessionCache.GetRefreshTokenId(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		return err
	}

	refreshSession, err := u.SessionCache.GetSession(ctx, refreshTokenId)
	if err != nil {
		return err
	}

	if refreshSession.TokenHash != hashRefreshToken(refreshToken) {
		return fmt.Errorf("невозможно выполнить операцию")
	}

	if refreshSession.IP != clientIP {
		return fmt.Errorf("невозможно выполнить операцию")
	}

	if refreshSession.UserAgent != ua {
		return fmt.Errorf("невозможно выполнить операцию")
	}

	err = u.SessionCache.DeleteAllUserSessions(ctx, refreshSession.UserID)
	if err != nil {
		return err
	}

	return nil
}

func hashRefreshToken(token string) string {
	return hasher.Sha256Hex(token)
}
