package repository

import (
	"context"
	"full-project-mock/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, email string) (*model.User, error)
	Exists(ctx context.Context, email string) (bool, error)
}
