package service

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"full-project-mock/internal/domain/model"
	"full-project-mock/internal/domain/usecase"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type TokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	accessTTL  time.Duration
}

func NewTokenService(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey, ttl time.Duration) usecase.TokenService {
	return &TokenService{
		privateKey: privateKey,
		publicKey:  publicKey,
		accessTTL:  ttl,
	}
}

func (s *TokenService) GenerateAccessToken(user *model.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *TokenService) GenerateRefreshToken() (tokenID, plainToken string, err error) {
	tokenID = uuid.NewString()
	plainTokenBytes := make([]byte, 32)
	if _, err = rand.Read(plainTokenBytes); err != nil {
		return
	}

	plainToken = base64.RawStdEncoding.EncodeToString(plainTokenBytes)
	return
}

func (s *TokenService) ParseToken(tokenStr string) (jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return jwt.RegisteredClaims{}, fmt.Errorf("token is invalid")
	}

	if claims.ExpiresAt == nil {
		return jwt.RegisteredClaims{}, errors.New("expired at missing")
	}

	return *claims, nil
}
