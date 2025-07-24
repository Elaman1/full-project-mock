package user

import (
	"context"
	"fmt"
	domcache "github.com/Elaman1/full-project-mock/internal/domain/cache"
	"github.com/Elaman1/full-project-mock/internal/domain/constants"
	"github.com/Elaman1/full-project-mock/internal/domain/model"
	"github.com/Elaman1/full-project-mock/internal/domain/repository"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"github.com/Elaman1/full-project-mock/pkg/hasher"
	"net/http"
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

func (u *Usecase) Login(ctx context.Context, email, password, clientIP, ua string) (string, string, int, error) {
	user, err := u.Rep.Get(ctx, email)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	err = hasher.Verify(user.Password, password)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	return u.generateAccessAndRefreshToken(ctx, clientIP, ua, user)
}

func (u *Usecase) generateAccessAndRefreshToken(ctx context.Context, clientIP, ua string, user *model.User) (string, string, int, error) {
	accessToken, err := u.TokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	refreshTokenId, plainToken, err := u.TokenService.GenerateRefreshToken()
	if err != nil {
		return "", "", http.StatusInternalServerError, err
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
		return "", "", http.StatusInternalServerError, err
	}

	err = u.SessionCache.SaveSession(ctx, newSess, u.RefreshTtl)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	return accessToken, plainToken, http.StatusOK, nil
}

func (u *Usecase) Refresh(ctx context.Context, accessToken, refreshToken, clientIP, ua string) (string, string, int, error) {
	mapClaims, err := u.TokenService.ParseToken(accessToken)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	userId, err := strconv.Atoi(mapClaims.Subject)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	user, err := u.Rep.GetById(ctx, int64(userId))
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	refreshTokenId, err := u.SessionCache.GetRefreshTokenId(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	refreshSession, err := u.SessionCache.GetSession(ctx, refreshTokenId)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	if refreshSession.ExpiresAt.Before(time.Now()) {
		return "", "", http.StatusUnauthorized, fmt.Errorf("время истек заново авторизуйтесь")
	}

	if refreshSession.IP != clientIP {
		return "", "", http.StatusUnauthorized, fmt.Errorf("авторизуйтесь еще раз")
	}

	if refreshSession.UserAgent != ua {
		return "", "", http.StatusUnauthorized, fmt.Errorf("авторизуйтесь еще раз")
	}

	if hashRefreshToken(refreshToken) != refreshSession.TokenHash {
		return "", "", http.StatusUnauthorized, fmt.Errorf("авторизуйтесь еще раз")
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
