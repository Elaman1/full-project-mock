package user

import (
	"context"
	"github.com/Elaman1/full-project-mock/internal/middleware"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMeHandler(t *testing.T) {
	handler := &UserHandler{}

	tests := []struct {
		name         string
		ctx          context.Context
		expectedCode int
		expectedBody string
	}{
		{
			name: "User ID found in context",
			ctx: func() context.Context {
				ctx := context.Background()
				return middleware.SetUserIDToContext(ctx, "123")
			}(),
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"123"}`,
		},
		{
			name:         "User ID not found in context",
			ctx:          context.Background(),
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"Not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			req = req.WithContext(tt.ctx)

			rec := httptest.NewRecorder()
			handler.MeHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}
