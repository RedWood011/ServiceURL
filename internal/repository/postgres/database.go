package postgres

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Repository Repo.....
type Repository struct {
	DB *pgxpool.Pool
}

// NewDatabase Создать бд.....
func NewDatabase(ctx context.Context, dsn string, maxAttempts int) (db *Repository, err error) {
	var pool *pgxpool.Pool

	err = doWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	ok, err := startMigration(dsn)
	if err != nil && !ok {
		return nil, fmt.Errorf("failed migrate database: %w", err)
	}

	return &Repository{DB: pool}, err
}

func startMigration(dsn string) (bool, error) {
	_, b, _, _ := runtime.Caller(-0)
	basePath := filepath.Dir(b)
	migrationsPath := basePath + "/migrations"
	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}

	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}
	return true, nil
}

// Ping Пинг...
func (r *Repository) Ping(ctx context.Context) error {
	return r.DB.Ping(ctx)
}

// SaveDone Сохранение...
func (r *Repository) SaveDone() error {
	r.DB.Close()
	return nil
}
