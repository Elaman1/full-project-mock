package middleware

import (
	"errors"
	"full-project-mock/internal/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	testUserID = "123"
)

func TestAuthMiddleware(t *testing.T) {
	type testCase struct {
		name              string
		authHeader        string
		mockParseResponse jwt.RegisteredClaims
		mockParseErr      error
		expectedStatus    int
		expectedBody      string
		expectNextCalled  bool
	}

	now := time.Now()
	validClaims := jwt.RegisteredClaims{
		Subject:   testUserID,
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
	}
	expiredClaims := jwt.RegisteredClaims{
		Subject:   testUserID,
		ExpiresAt: jwt.NewNumericDate(now.Add(-10 * time.Minute)),
	}

	tests := []testCase{
		{
			name:             "missing authorization header",
			authHeader:       "",
			expectedStatus:   http.StatusUnauthorized,
			expectedBody:     "missing or invalid Authorization header",
			expectNextCalled: false,
		},
		{
			name:             "invalid authorization format",
			authHeader:       "InvalidFormat token",
			expectedStatus:   http.StatusUnauthorized,
			expectedBody:     "missing or invalid Authorization header",
			expectNextCalled: false,
		},
		{
			name:             "invalid token",
			authHeader:       "Bearer bad-token",
			mockParseErr:     errors.New("invalid token"),
			expectedStatus:   http.StatusUnauthorized,
			expectedBody:     "invalid token",
			expectNextCalled: false,
		},
		{
			name:              "expired token",
			authHeader:        "Bearer expired-token",
			mockParseResponse: expiredClaims,
			expectedStatus:    http.StatusUnauthorized,
			expectedBody:      "expired token",
			expectNextCalled:  false,
		},
		{
			name:              "empty subject",
			authHeader:        "Bearer empty-subject",
			mockParseResponse: jwt.RegisteredClaims{Subject: "", ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute))},
			expectedStatus:    http.StatusUnauthorized,
			expectedBody:      "invalid token: empty subject",
			expectNextCalled:  false,
		},
		{
			name:              "valid token",
			authHeader:        "Bearer good-token",
			mockParseResponse: validClaims,
			expectedStatus:    http.StatusOK,
			expectedBody:      "",
			expectNextCalled:  true,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockTokenSvc := new(mocks.MockTokenService)

			// Мокаем ParseToken, если передан header
			if tc.authHeader != "" && strings.HasPrefix(tc.authHeader, "Bearer ") {
				tokenStr := strings.TrimPrefix(tc.authHeader, "Bearer ")
				mockTokenSvc.On("ParseToken", tokenStr).
					Return(tc.mockParseResponse, tc.mockParseErr).Maybe()
			}

			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true

				uid, ok := GetUserIDFromContext(r.Context())
				require.True(t, ok)
				require.Equal(t, testUserID, uid)

				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()

			middlewareFunc := AuthMiddleware(mockTokenSvc)
			handler := middlewareFunc(nextHandler)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectNextCalled, nextCalled)

			// Проверка тела ответа (кроме валидного запроса)
			if tc.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}

			mockTokenSvc.AssertExpectations(t)
		})
	}
}
