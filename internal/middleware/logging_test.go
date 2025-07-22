package middleware

import (
	"bytes"
	"full-project-mock/internal/service"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogMiddleware(t *testing.T) {
	type testCase struct {
		name           string
		statusToReturn int
		expectLog      bool
	}

	tests := []testCase{
		{
			name:           "status OK should log",
			statusToReturn: http.StatusOK,
			expectLog:      true,
		},
		{
			name:           "status NotFound should not log",
			statusToReturn: http.StatusNotFound,
			expectLog:      false,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Ловим вывод логгера
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo}))

			handlerCalled := false
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true

				// Проверим, что логгер есть в контексте
				logFromCtx := service.LoggerFromContext(r.Context())
				assert.NotNil(t, logFromCtx)

				w.WriteHeader(tc.statusToReturn)
			})

			req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
			rec := httptest.NewRecorder()

			middleware := LogMiddleware(logger)
			middleware(handler).ServeHTTP(rec, req)

			assert.True(t, handlerCalled)
			assert.Equal(t, tc.statusToReturn, rec.Code)

			output := logBuf.String()
			if tc.expectLog {
				assert.Contains(t, output, "request completed")
				assert.Contains(t, output, "path=test-path")
				assert.Contains(t, output, "method=GET")
			} else {
				assert.NotContains(t, output, "request completed")
			}
		})
	}
}
