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

	claims := jwt.MapClaims{
		"sub":  strconv.Itoa(int(user.ID)),         // subject = user ID
		"exp":  time.Now().Add(s.accessTTL).Unix(), // expiration
		"iat":  time.Now().Unix(),                  // issued at
		"role": user.Role,                          // возможно
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

func (s *TokenService) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("signature invalid")
	}

	// проверка exp вручную
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("exp claim missing or wrong type")
	}

	if int64(exp) < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
