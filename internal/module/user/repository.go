package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Elaman1/full-project-mock/internal/domain/model"
	"github.com/Elaman1/full-project-mock/internal/domain/repository"
	"time"
)

type Repository struct {
	DB DBExecutor
}

// DBExecutor Чтобы можно было или sql.DB || sql.Tx передавать
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewUserRepository(db DBExecutor) repository.UserRepository {
	return &Repository{
		DB: db,
	}
}

func (u *Repository) Create(ctx context.Context, user *model.User) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := u.DB.ExecContext(ctxTimeout, "INSERT INTO users (email, password, name, role_id) VALUES ($1, $2, $3, $4)", user.Email, user.Password, user.Username, user.RoleID)
	if err != nil {
		return fmt.Errorf("create user error: %w", err)
	}

	return nil
}

func (u *Repository) Get(ctx context.Context, email string) (*model.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user := &model.User{}
	err := u.DB.QueryRowContext(ctxTimeout, "SELECT id, email, name, password, created_at, role_id FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.RoleID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}

func (u *Repository) Exists(ctx context.Context, email string) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var dummy int
	err := u.DB.QueryRowContext(ctxTimeout, "select 1 from users where email = $1 limit 1", email).Scan(&dummy)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, err
}

func (u *Repository) GetById(ctx context.Context, id int64) (*model.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	user := &model.User{}
	err := u.DB.QueryRowContext(ctxTimeout, "SELECT id, email, name, password, created_at, role_id FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.RoleID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}
