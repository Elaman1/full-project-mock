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
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var (
	testDB    *sql.DB
	testRedis *redis.Client
	accessTTL time.Duration
)

func TestMain(m *testing.M) {
	err := initConn("../../../config/config.test.yaml")
	if err != nil {
		panic(err)
	}

	code := m.Run()
	err = testDB.Close()
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

func initConn(path string) error {
	confFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg config.Config
	if err = yaml.Unmarshal(confFile, &cfg); err != nil {
		return err
	}

	err = validateCfg(cfg)
	if err != nil {
		return err
	}

	conn, err := initTestDB(&cfg)
	if err != nil {
		return err
	}

	r, err := initRedis(&cfg)
	if err != nil {
		return err
	}

	ttl, err := time.ParseDuration(cfg.JWT.AccessTTL)
	if err != nil {
		return err
	}

	accessTTL = ttl
	testDB = conn
	testRedis = r
	return nil
}

func validateCfg(cfg config.Config) error {
	if err := validatePostgresDB(&cfg); err != nil {
		return err
	}

	if err := validateRedis(&cfg); err != nil {
		return err
	}

	if err := validateJWTAccessTTL(&cfg); err != nil {
		return err
	}

	return nil
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

func validateRedis(cfg *config.Config) error {
	if cfg.Redis.Host == "" {
		return errors.New("missing required test configuration variable: redis_host")
	}

	if cfg.Redis.Port == 0 {
		return errors.New("missing required test configuration variable: redis_port")
	}

	return nil
}

func validateJWTAccessTTL(cfg *config.Config) error {
	if cfg.JWT.AccessTTL == "" {
		return errors.New("missing required test configuration variable: jwt_access_ttl")
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

func initTestDB(cfg *config.Config) (*sql.DB, error) {
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

func initRedis(cfg *config.Config) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	status := r.Ping(ctx)
	if err := status.Err(); err != nil {
		return nil, err
	}

	return r, nil
}
