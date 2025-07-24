package user

import (
	"context"
	"encoding/json"
	"github.com/Elaman1/full-project-mock/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLoginHandler_Integration(t *testing.T) {
	tx := initTx(t)

	// Регистрация, чтобы проверить логин
	registerHand, _, _ := buildUserHandlerIntegration(t, tx)
	registerReq := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"login-user@test.com","username":"integrationUser","password":"secure123"}`))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq = registerReq.WithContext(service.WithLogger(registerReq.Context(), slog.Default()))
	registerRec := httptest.NewRecorder()
	registerHand.RegisterHandler(registerRec, registerReq)

	registerRes := registerRec.Result()
	defer registerRes.Body.Close()
	require.Equal(t, http.StatusCreated, registerRes.StatusCode)

	body, err := io.ReadAll(registerRes.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Успешно создано")

	// После успешной регистрации проверяем логин
	tests := []struct {
		name           string
		payload        string
		wantStatus     int
		wantInResponse string
		checkResponse  bool
	}{
		{
			name:           "success",
			payload:        `{"email":"login-user@test.com","password":"secure123"}`,
			wantStatus:     http.StatusOK,
			wantInResponse: "access_token",
			checkResponse:  true,
		},
		{
			name:           "empty_request",
			payload:        ``,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: `Invalid request payload`,
		},
		{
			name:           "error_request",
			payload:        `{"bad": `,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: `Invalid request payload`,
		},
		{
			name:           "empty_request",
			payload:        `{"email":"","password":""}`,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: "email or password is empty",
		},
		{
			name:           "incorrect_email",
			payload:        `{"email":"bad-email","password":"not-correct-password"}`,
			wantStatus:     http.StatusBadRequest,
			wantInResponse: "invalid email",
		},
		{
			name:           "not_found_email",
			payload:        `{"email":"not-email@test.com","password":"not-correct-password"}`,
			wantStatus:     http.StatusUnauthorized,
			wantInResponse: "user not found",
		},
		{
			name:           "not_correct_password",
			payload:        `{"email":"login-user@test.com","password":"not"}`,
			wantStatus:     http.StatusUnauthorized,
			wantInResponse: "логин или пароль неправильный",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, tokenSrv, redisRepo := buildUserHandlerIntegration(t, tx)
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.payload))

			req = req.WithContext(service.WithLogger(req.Context(), slog.Default()))
			req.Header.Set("User-Agent", testAgent)
			rec := httptest.NewRecorder()

			handler.LoginHandler(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			loginBody, loginErr := io.ReadAll(res.Body)

			if tt.checkResponse {
				var respStruct LoginAndRefreshResponseStruct
				respErr := json.Unmarshal(loginBody, &respStruct)
				if respErr != nil {
					t.Fatalf("response error: %v", respErr)
				}

				// Проверка на пустоту
				assert.NotEmpty(t, respStruct.AccessToken)
				assert.NotEmpty(t, respStruct.RefreshToken)

				// Получаем данные с редиса, чтобы проверить и дальше проверки токенов
				hashedPlainToken := hashRefreshToken(respStruct.RefreshToken)
				tokenID, getIDErr := redisRepo.GetRefreshTokenId(context.Background(), hashedPlainToken)
				require.NoError(t, getIDErr)
				require.NotEmpty(t, tokenID)

				refreshSession, getErr := redisRepo.GetSession(context.Background(), tokenID)
				require.NoError(t, getErr)
				require.NotNil(t, refreshSession)
				require.NotZero(t, refreshSession.UserID)
				require.Equal(t, testAgent, refreshSession.UserAgent)

				// Проверяем токены
				claims, parseErr := tokenSrv.ParseToken(respStruct.AccessToken)
				assert.NoError(t, parseErr)
				assert.Equal(t, strconv.Itoa(int(refreshSession.UserID)), claims.Subject)
				assert.True(t, claims.ExpiresAt.After(time.Now()))
			}

			assert.NoError(t, loginErr)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, string(loginBody), tt.wantInResponse, "unexpected body for test: %s", string(loginBody))
		})
	}
}
