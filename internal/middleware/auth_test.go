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

// Вспомогательный handler, который проверяет наличие userID в контексте
func newTestHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := GetUserIDFromContext(r.Context())
		require.True(t, ok)
		require.Equal(t, testUserID, uid)
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuthMiddleware(t *testing.T) {
	type testCase struct {
		name              string
		authHeader        string
		mockParseResponse jwt.RegisteredClaims
		mockParseErr      error
		expectedStatus    int
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
			expectNextCalled: false,
		},
		{
			name:             "invalid authorization format",
			authHeader:       "InvalidFormat token",
			expectedStatus:   http.StatusUnauthorized,
			expectNextCalled: false,
		},
		{
			name:             "invalid token",
			authHeader:       "Bearer bad-token",
			mockParseErr:     errors.New("invalid token"),
			expectedStatus:   http.StatusUnauthorized,
			expectNextCalled: false,
		},
		{
			name:              "expired token",
			authHeader:        "Bearer expired-token",
			mockParseResponse: expiredClaims,
			expectedStatus:    http.StatusUnauthorized,
			expectNextCalled:  false,
		},
		{
			name:              "valid token",
			authHeader:        "Bearer good-token",
			mockParseResponse: validClaims,
			expectedStatus:    http.StatusOK,
			expectNextCalled:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTokenSvc := new(mocks.MockTokenService)

			// если ожидается вызов ParseToken — настраиваем мок
			if tc.authHeader != "" && strings.HasPrefix(tc.authHeader, "Bearer ") {
				tokenStr := strings.TrimPrefix(tc.authHeader, "Bearer ")
				mockTokenSvc.On("ParseToken", tokenStr).
					Return(tc.mockParseResponse, tc.mockParseErr).Maybe()
			}

			// фиксируем был ли вызван следующий handler
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				newTestHandler(t).ServeHTTP(w, r)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()

			AuthMiddleware(mockTokenSvc)(nextHandler).ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectNextCalled, nextCalled)
			mockTokenSvc.AssertExpectations(t)
		})
	}
}
