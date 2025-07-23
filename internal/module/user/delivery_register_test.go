package user

import (
	"context"
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

func TestRegisterHandler(t *testing.T) {
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
			name: "Invalid JSON body",
			args: args{
				body:         "invalid-json",
				mockSetup:    func(m *MockUserUsecase) {},
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"error":"Invalid request payload"}`,
			},
		},
		{
			name: "Validation error",
			args: args{
				body:         `{"email":"bad","username":"","password":""}`,
				mockSetup:    func(m *MockUserUsecase) {},
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"error":"User validation error: `,
			},
		},
		{
			name: "Usecase returns error",
			args: args{
				body: fmt.Sprintf(`{"email":"%s","username":"user","password":"pass123"}`, email),
				mockSetup: func(m *MockUserUsecase) {
					m.On("Register", mock.Anything, "test@mail.com", "user", "pass123").
						Return(int64(0), errors.New("some error")).Once()
				},
				expectedCode: http.StatusInternalServerError,
				expectedBody: `{"error":"User registration error: some error"}`,
			},
		},
		{
			name: "Success",
			args: args{
				body: fmt.Sprintf(`{"email":"%s","username":"user","password":"pass123"}`, email),
				mockSetup: func(m *MockUserUsecase) {
					m.On("Register", mock.Anything, email, "user", "pass123").
						Return(int64(123), nil).Once()
				},
				expectedCode: http.StatusCreated,
				expectedBody: `{"success":"Успешно создано"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockUserUsecase)
			tt.args.mockSetup(mockUsecase)

			handler := &UserHandler{Usecase: mockUsecase}

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.args.body))
			req = req.WithContext(service.WithLogger(req.Context(), slog.Default()))

			rec := httptest.NewRecorder()

			handler.RegisterHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.args.expectedCode, res.StatusCode)
			assert.Contains(t, string(body), tt.args.expectedBody)

			mockUsecase.AssertExpectations(t)
		})
	}
}

type MockUserUsecase struct {
	mock.Mock
}

func (mock *MockUserUsecase) Register(ctx context.Context, email, username, password string) (int64, error) {
	args := mock.Called(ctx, email, username, password)
	id, ok := args.Get(0).(int64)
	if !ok {
		return 0, args.Error(1)
	}

	return id, args.Error(1)
}

func (mock *MockUserUsecase) Login(ctx context.Context, email, password, clientIP, ua string) (string, string, error) {
	args := mock.Called(ctx, email, password, clientIP, ua)
	return args.String(0), args.String(1), args.Error(2)
}

func (mock *MockUserUsecase) Logout(ctx context.Context, refreshToken, clientIP, ua string) error {
	args := mock.Called(ctx, refreshToken, clientIP, ua)
	return args.Error(0)
}

func (mock *MockUserUsecase) LogoutAllDevices(ctx context.Context, refreshToken, clientIP, ua string) error {
	args := mock.Called(ctx, refreshToken, clientIP, ua)
	return args.Error(0)
}

func (mock *MockUserUsecase) Refresh(ctx context.Context, accessToken, refreshToken, clientIP, ua string) (string, string, error) {
	args := mock.Called(ctx, accessToken, refreshToken, clientIP, ua)
	return args.String(0), args.String(1), args.Error(2)
}
