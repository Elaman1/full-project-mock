package user

import (
	"context"
	"errors"
	"fmt"
	"full-project-mock/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	defaultEmail    = "test@example.com"
	defaultUserName = "testuser"
	defaultPassword = "securepassword"
)

func TestRegister_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)

	mockRepo.On("Exists", ctx, defaultEmail).Return(false, nil)
	mockRepo.On("Create", ctx, mock.Anything).Return(nil)

	uc := &Usecase{
		Rep: mockRepo,
	}

	id, err := uc.Register(ctx, defaultEmail, defaultUserName, defaultPassword)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), id)

	mockRepo.AssertExpectations(t)
}

func TestRegister_Exists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)

	mockRepo.On("Exists", ctx, defaultEmail).Return(true, nil)

	uc := &Usecase{
		Rep: mockRepo,
	}

	_, err := uc.Register(ctx, defaultEmail, defaultUserName, defaultPassword)

	assert.EqualError(t, err, fmt.Sprintf("пользователь с таким email %s уже существует", defaultEmail))

	mockRepo.AssertExpectations(t)
}

func TestRegister_ExistsError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)

	customError := errors.New("custom error")
	mockRepo.On("Exists", ctx, defaultEmail).Return(false, customError)

	uc := &Usecase{
		Rep: mockRepo,
	}

	_, err := uc.Register(ctx, defaultEmail, defaultUserName, defaultPassword)

	assert.EqualError(t, err, customError.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_CreateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)

	customerError := errors.New("customer error")

	mockRepo.On("Exists", ctx, defaultEmail).Return(false, nil)
	mockRepo.On("Create", ctx, mock.Anything).Return(customerError)

	uc := &Usecase{
		Rep: mockRepo,
	}

	_, err := uc.Register(ctx, defaultEmail, defaultUserName, defaultPassword)

	assert.EqualError(t, err, "произошла ошибка при регистрации")
	mockRepo.AssertExpectations(t)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Get(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}

	user, ok := u.(*model.User)
	if !ok {
		return nil, fmt.Errorf("error casting model.User")
	}

	return user, args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	user, ok := args.Get(0).(*model.User)
	if !ok {
		return nil, fmt.Errorf("error casting model.User")
	}

	return user, args.Error(1)
}
