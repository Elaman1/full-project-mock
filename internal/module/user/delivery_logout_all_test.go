package user

import (
	"errors"
	"fmt"
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

var (
	refreshStr = "refresh123"
)

func TestLogoutAllHandler(t *testing.T) {
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
				body:         `{"refresh_token":"abc`, // malformed
				mockSetup:    func(m *MockUserUsecase) {},
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"error":"Invalid request payload"}`,
			},
		},
		{
			name: "Usecase returns error",
			args: args{
				body: fmt.Sprintf(`{"refresh_token":"%s"}`, refreshStr),
				mockSetup: func(m *MockUserUsecase) {
					m.On("LogoutAllDevices", mock.Anything, refreshStr, ipAddress, testAgent).
						Return(errors.New("something went wrong")).Once()
				},
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"error":"error logout all"}`,
			},
		},
		{
			name: "Success",
			args: args{
				body: fmt.Sprintf(`{"refresh_token":"%s"}`, refreshStr),
				mockSetup: func(m *MockUserUsecase) {
					m.On("LogoutAllDevices", mock.Anything, refreshStr, ipAddress, testAgent).
						Return(nil).Once()
				},
				expectedCode: http.StatusOK,
				expectedBody: ``, // Тело пустое — это ожидаемо
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.args.mockSetup(mockUsecase)

			handler := &UserHandler{Usecase: mockUsecase}

			req := httptest.NewRequest(http.MethodPost, "/logout-all", strings.NewReader(tt.args.body))
			req.Header.Set("User-Agent", testAgent)
			req.RemoteAddr = ipAddress
			req = req.WithContext(service.WithLogger(req.Context(), slog.Default()))

			rec := httptest.NewRecorder()
			handler.LogoutAllHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.args.expectedCode, res.StatusCode)

			if tt.args.expectedBody != "" {
				assert.Contains(t, string(body), tt.args.expectedBody)
			} else {
				assert.Empty(t, string(body)) // тело пустое при успехе
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
