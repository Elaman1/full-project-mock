package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"full-project-mock/internal/config"
	"full-project-mock/internal/domain/constants"
	"full-project-mock/internal/domain/model"
	"full-project-mock/internal/domain/repository"
	"full-project-mock/pkg/hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testDB *sql.DB
)

func TestMain(m *testing.M) {
	db, err := initDB("config/config.test.yaml")
	if err != nil {
		panic(err)
	}

	testDB = db

	code := m.Run()
	err = db.Close()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	user, repo, err := createTestUser(t, ctx, 1)
	require.NoError(t, err)

	u, err := repo.Get(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Email, u.Email)
	require.Equal(t, user.Username, u.Username)
	require.Equal(t, user.RoleID, u.RoleID)
	require.Equal(t, user.Password, u.Password)
}

func TestExists(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	user, repo, err := createTestUser(t, ctx, 2)
	require.NoError(t, err)
	ok, err := repo.Exists(ctx, user.Email)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestGetById(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	user, repo, err := createTestUser(t, ctx, 3)
	require.NoError(t, err)
	savedUser, err := repo.Get(ctx, user.Email)
	require.NoError(t, err)

	foundedUser, err := repo.GetById(ctx, savedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, savedUser.ID, foundedUser.ID)
	assert.Equal(t, savedUser.Email, foundedUser.Email)
}

func createTestUser(t *testing.T, ctx context.Context, prefix int) (*model.User, repository.UserRepository, error) {
	t.Helper()

	userEmail := fmt.Sprintf("%s-%d", defaultEmail, prefix) // Чтобы унифицировать
	pwd, err := hasher.HashPassword(defaultPassword)
	require.NoError(t, err)

	db := setupTestDB(t)
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewUserRepository(tx)
	t.Cleanup(func() {
		// Не будем хранить данные, а сразу откатывать
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			t.Fatal(rollbackErr)
		}
	})

	user := &model.User{
		Email:    userEmail,
		Password: pwd,
		Username: defaultUserName,
		RoleID:   constants.DefaultUserRoleID,
	}
	err = repo.Create(ctx, user)
	return user, repo, err
}

func initDB(path string) (*sql.DB, error) {
	confFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	if err = yaml.Unmarshal(confFile, &cfg); err != nil {
		return nil, err
	}

	if err = validatePostgresDB(&cfg); err != nil {
		return nil, err
	}

	c := cfg.PostgresDB
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SslMode,
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

func validatePostgresDB(cfg *config.Config) error {
	if cfg.PostgresDB.Host == "" {
		return errors.New("missing required test configuration variable: postgres_host")
	}

	if cfg.PostgresDB.Port == 0 {
		return errors.New("missing required test configuration variable: postgres_port")
	}

	if cfg.PostgresDB.DBName == "" {
		return errors.New("missing required test configuration variable: postgres_dbname")
	}

	if cfg.PostgresDB.SslMode == "" {
		return errors.New("missing required test configuration variable: postgres_ssl_mode")
	}

	return nil
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}
	return testDB
}
