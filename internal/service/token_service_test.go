package service

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/Elaman1/full-project-mock/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return privateKey, &privateKey.PublicKey
}

func TestTokenService_ParseToken(t *testing.T) {
	privateKey, publicKey := generateTestKeys(t)
	tokenSvc := NewTokenService(publicKey, privateKey, time.Minute)

	fixedNow := time.Now()

	// Валидный токен
	validToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   "valid-subject",
		ExpiresAt: jwt.NewNumericDate(fixedNow.Add(10 * time.Minute)),
	})
	validSigned, err := validToken.SignedString(privateKey)
	require.NoError(t, err)

	// Просроченный токен
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   "expired-subject",
		ExpiresAt: jwt.NewNumericDate(fixedNow.Add(-10 * time.Minute)),
	})
	expiredSigned, err := expiredToken.SignedString(privateKey)
	require.NoError(t, err)

	// Подпись с другим ключом
	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	wrongSigToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   "bad-signature",
		ExpiresAt: jwt.NewNumericDate(fixedNow.Add(10 * time.Minute)),
	})
	wrongSigSigned, err := wrongSigToken.SignedString(otherKey)
	require.NoError(t, err)

	// Подпись с неподдерживаемым алгоритмом (HS256)
	badAlgToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "bad-alg",
		ExpiresAt: jwt.NewNumericDate(fixedNow.Add(10 * time.Minute)),
	})
	badAlgSigned, err := badAlgToken.SignedString([]byte("secret"))
	require.NoError(t, err)

	// Токен с валидной подписью, но token.Valid == false (невалидный claims)
	noExpToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject: "no-exp", // отсутствует ExpiresAt
	})
	noExpSigned, err := noExpToken.SignedString(privateKey)
	require.NoError(t, err)

	tests := []struct {
		name         string
		inputToken   string
		wantSubject  string
		wantExpired  bool
		expectErr    bool
		expectErrMsg string
	}{
		{
			name:        "valid token",
			inputToken:  validSigned,
			wantSubject: "valid-subject",
			wantExpired: false,
			expectErr:   false,
		},
		{
			name:        "expired token",
			inputToken:  expiredSigned,
			wantSubject: "expired-subject",
			wantExpired: true,
			expectErr:   true,
		},
		{
			name:         "malformed token",
			inputToken:   "not.a.jwt",
			expectErr:    true,
			expectErrMsg: "token is malformed",
		},
		{
			name:         "wrong signature",
			inputToken:   wrongSigSigned,
			expectErr:    true,
			expectErrMsg: "failed to parse token",
		},
		{
			name:         "unexpected signing algorithm",
			inputToken:   badAlgSigned,
			expectErr:    true,
			expectErrMsg: "unexpected signing method",
		},
		{
			name:         "token invalid (token.Valid == false)",
			inputToken:   noExpSigned,
			expectErr:    true,
			expectErrMsg: "expired at missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, parseErr := tokenSvc.ParseToken(tt.inputToken)

			if tt.expectErr {
				require.Error(t, parseErr)
				if tt.expectErrMsg != "" {
					assert.Contains(t, parseErr.Error(), tt.expectErrMsg)
				}
				return
			}

			require.NoError(t, parseErr)
			assert.Equal(t, tt.wantSubject, claims.Subject)

			// Безопасная проверка на истечение
			if claims.ExpiresAt != nil {
				expired := claims.ExpiresAt.Time.Before(fixedNow)
				assert.Equal(t, tt.wantExpired, expired)
			} else {
				assert.False(t, tt.wantExpired, "ExpiresAt is nil, but wantExpired is true")
			}
		})
	}
}

func TestTokenService_GenerateAccessToken_NilUser(t *testing.T) {
	_, publicKey := generateTestKeys(t)
	svc := NewTokenService(publicKey, nil, time.Minute)

	token, err := svc.GenerateAccessToken(nil)
	assert.Empty(t, token)
	assert.Error(t, err)
	assert.EqualError(t, err, "user is nil")
}

func TestTokenService_GenerateAccessToken(t *testing.T) {
	privateKey, publicKey := generateTestKeys(t)
	tokenSvc := NewTokenService(publicKey, privateKey, time.Minute)

	now := time.Now()

	tests := []struct {
		name        string
		user        *model.User
		expectErr   bool
		expectEmpty bool
	}{
		{
			name:        "nil user",
			user:        nil,
			expectErr:   true,
			expectEmpty: true,
		},
		{
			name: "valid user",
			user: &model.User{
				ID: 42,
			},
			expectErr:   false,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr, err := tokenSvc.GenerateAccessToken(tt.user)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, tokenStr)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, tokenStr)

			// Парсим токен обратно и проверяем claims
			parsedToken, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			require.NoError(t, err)
			require.True(t, parsedToken.Valid)

			claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
			require.True(t, ok)

			assert.Equal(t, "42", claims.Subject)

			// Допуск ±2 секунды
			expectedExp := now.Add(time.Minute)
			diff := claims.ExpiresAt.Time.Sub(expectedExp)
			assert.LessOrEqual(t, diff.Abs(), 2*time.Second)
		})
	}
}

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	privateKey, publicKey := generateTestKeys(t)
	tokenSvc := NewTokenService(publicKey, privateKey, time.Minute)

	tokenID, plainToken, err := tokenSvc.GenerateRefreshToken()

	require.NoError(t, err)
	assert.NotEmpty(t, tokenID)
	assert.NotEmpty(t, plainToken)

	// Проверка UUID
	_, uuidErr := uuid.Parse(tokenID)
	assert.NoError(t, uuidErr)

	// Проверка base64 длины (32 байта → 43-44 символа без паддинга)
	assert.GreaterOrEqual(t, len(plainToken), 43)
}
