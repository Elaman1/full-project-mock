package user

import (
	"errors"
	"fmt"
	"github.com/Elaman1/full-project-mock/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	ipAddress   = "127.0.0.1"
	email       = "test@mail.com"
	wrongPass   = "wrongpass"
	correctPass = "correctpass"
	testAgent   = "test-agent"
)

func TestLoginHandler(t *testing.T) {
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
				body:         fmt.Sprintf(`{"email": "%s",`, email), // malformed JSON
				mockSetup:    func(m *MockUserUsecase) {},
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"error":"Invalid request payload"}`,
			},
		},
		{
			name: "Usecase returns error",
			args: args{
				body: fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, wrongPass),
				mockSetup: func(m *MockUserUsecase) {
					m.On("Login", mock.Anything, email, wrongPass, ipAddress, testAgent).
						Return("", "", http.StatusBadRequest, errors.New("invalid credentials")).Once()
				},
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"error":"User login error: invalid credentials"}`,
			},
		},
		{
			name: "Success",
			args: args{
				body: fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, correctPass),
				mockSetup: func(m *MockUserUsecase) {
					m.On("Login", mock.Anything, email, correctPass, ipAddress, testAgent).
						Return("access-token", "refresh-token", http.StatusOK, nil).Once()
				},
				expectedCode: http.StatusOK,
				expectedBody: `"access_token":"access-token"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.args.mockSetup(mockUsecase)

			handler := &UserHandler{Usecase: mockUsecase}

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.args.body))
			req.Header.Set("User-Agent", testAgent)
			req.RemoteAddr = ipAddress
			req = req.WithContext(service.WithLogger(req.Context(), slog.Default()))

			rec := httptest.NewRecorder()

			handler.LoginHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.args.expectedCode, res.StatusCode)
			assert.Contains(t, string(body), tt.args.expectedBody)

			mockUsecase.AssertExpectations(t)
		})
	}
}
