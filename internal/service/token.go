package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"full-project-mock/internal/domain/model"
	"full-project-mock/internal/domain/usecase"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type TokenService struct {
	secret    []byte
	accessTTL time.Duration
}

func NewTokenService(secret string, ttl time.Duration) usecase.TokenService {
	return &TokenService{
		secret:    []byte(secret),
		accessTTL: ttl,
	}
}

// GenerateAccessToken 1. Генерация access token (JWT)
func (s *TokenService) GenerateAccessToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(s.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateRefreshToken 2. Генерация refresh token (обычная строка)
func (s *TokenService) GenerateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// ParseToken 3. Разбор access token — достать userID (sub)
func (s *TokenService) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid sub")
		}
		return sub, nil
	}

	return "", errors.New("invalid token")
}
