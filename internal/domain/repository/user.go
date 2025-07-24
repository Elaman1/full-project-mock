package repository

import (
	"context"
	"github.com/Elaman1/full-project-mock/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, email string) (*model.User, error)
	Exists(ctx context.Context, email string) (bool, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
}
