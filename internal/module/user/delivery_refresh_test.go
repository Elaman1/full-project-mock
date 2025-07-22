package user

import (
	"errors"
	"full-project-mock/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRefreshHandler(t *testing.T) {
	type args struct {
		body         string
		mockSetup    func(m *MockUserUsecase)
		expectedCode int
		expectedBody string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Invalid JSON",
			args: args{
				body:         `{"access_token": "abc",`, // malformed
				mockSetup:    func(m *MockUserUsecase) {},
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"error":"Invalid request payload"}`,
			},
		},
		{
			name: "Usecase returns error",
			args: args{
				body: `{"access_token":"access123","refresh_token":"refresh123"}`,
				mockSetup: func(m *MockUserUsecase) {
					m.On("Refresh", mock.Anything, "access123", "refresh123", "127.0.0.1", "test-agent").
						Return("", "", errors.New("invalid token")).Once()
				},
				expectedCode: http.StatusInternalServerError,
				expectedBody: `{"error":"Refresh error: invalid token"}`,
			},
		},
		{
			name: "Success",
			args: args{
				body: `{"access_token":"access123","refresh_token":"refresh123"}`,
				mockSetup: func(m *MockUserUsecase) {
					m.On("Refresh", mock.Anything, "access123", "refresh123", "127.0.0.1", "test-agent").
						Return("new-access", "new-refresh", nil).Once()
				},
				expectedCode: http.StatusCreated,
				expectedBody: `"access_token":"new-access"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.args.mockSetup(mockUsecase)

			handler := &UserHandler{Usecase: mockUsecase}

			req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(tt.args.body))
			req.Header.Set("User-Agent", "test-agent")
			req.RemoteAddr = "127.0.0.1:12345"
			req = req.WithContext(service.WithLogger(req.Context(), slog.Default()))

			rec := httptest.NewRecorder()
			handler.RefreshHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.args.expectedCode, res.StatusCode)
			assert.Contains(t, string(body), tt.args.expectedBody)

			mockUsecase.AssertExpectations(t)
		})
	}
}
