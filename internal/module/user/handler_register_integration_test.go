package user

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"errors"
	"github.com/Elaman1/full-project-mock/internal/cache"
	cache2 "github.com/Elaman1/full-project-mock/internal/domain/cache"
	"github.com/Elaman1/full-project-mock/internal/domain/usecase"
	"github.com/Elaman1/full-project-mock/internal/service"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Делаем все в одной транзакции так как мы не комитим, а откатываем
func initTx(t *testing.T) *sql.Tx {
	setConn(t)

	t.Helper()
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		// Не будем хранить данные, а сразу откатывать
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			t.Fatal(rollbackErr)
		}

		if rollbackRedis := testRedis.FlushDB(context.Background()).Err(); rollbackRedis != nil {
			t.Fatal("failed to flush redis:", rollbackRedis)
		}
	})

	return tx
}

// Сборка handler → usecase → repository на реальной БД
func buildUserHandlerIntegration(t *testing.T, tx *sql.Tx) (*UserHandler, usecase.TokenService, cache2.SessionCache) {
	t.Helper()

	userRepo := NewUserRepository(tx)
	sessionCache := cache.NewSessionRedisRepository(testRedis)
	privateKey, publicKey := generateTestKeys(t)
	tokenService := service.NewTokenService(publicKey, privateKey, accessTTL)
	usecase := NewUserUsecase(userRepo, tokenService, sessionCache)
	return &UserHandler{Usecase: usecase}, tokenService, sessionCache
}
func TestRegisterHandler_Integration(t *testing.T) {
	tx := initTx(t)

	tests := []struct {
		name           string
		payload        string
		wantStatus     int
		wantInResponse string
	}{
		{
			name:           "valid registration",
			payload:        `{"email":"int-user1@test.com","username":"integrationUser","password":"secure123"}`,
			wantStatus:     http.StatusCreated,
			wantInResponse: "Успешно создано",
		},
		{
			name:           "missing email",
			payload:        `{"username":"user","password":"secure123"}`,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: "email",
		},
		{
			name:           "missing password",
			payload:        `{"email":"int-user2@test.com","username":"user"}`,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: "password",
		},
		{
			name:           "empty username",
			payload:        `{"email":"int-user3@test.com","username":"","password":"secure123"}`,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: "username",
		},
		{
			name:           "invalid email format",
			payload:        `{"email":"bad-email","username":"user","password":"secure123"}`,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: "email",
		},
		{
			name:           "duplicate email",
			payload:        `{"email":"int-user4@test.com","username":"user","password":"secure123"}`,
			wantStatus:     http.StatusCreated,
			wantInResponse: "Успешно создано", // первый вызов
		},
		{
			name:           "duplicate email (again)",
			payload:        `{"email":"int-user4@test.com","username":"user2","password":"anotherpass"}`,
			wantStatus:     http.StatusInternalServerError, // чтобы не делить ошибки по тексту
			wantInResponse: "пользователь с таким email",   // второй вызов
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Используем отдельный handler для каждого кейса (чтобы rollback и RedisFlush работали корректно)
			handler, _, _ := buildUserHandlerIntegration(t, tx)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(service.WithLogger(req.Context(), slog.Default()))

			rec := httptest.NewRecorder()

			handler.RegisterHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, res.StatusCode, "unexpected status for test #%d", i)
			assert.Contains(t, string(body), tt.wantInResponse, "unexpected body for test #%d: %s", i, string(body))
		})
	}
}

func setConn(t *testing.T) {
	if testDB != nil && testRedis != nil {
		return
	}

	err := initConn("../../../config/config.test.yaml")
	if err != nil {
		t.Fatal(err)
	}
}

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("failed to generate RSA key:", err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}

type LoginAndRefreshResponseStruct struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
