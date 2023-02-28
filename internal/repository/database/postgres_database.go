package database

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/slog"
)

type Db interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Rollback(ctx context.Context) error
}

type Repository struct {
	conn   *pgxpool.Pool
	logger *slog.Logger
}

func (db *Database) NewRepository() Repository {
	return Repository{
		conn: db.conn,
	}
}

type Database struct {
	conn *pgxpool.Pool
}

func (db *Database) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}
func (db *Database) Close() error {
	db.conn.Close()
	return nil
}

func NewRepo(dsn string, maxAttempts int) (r *Repository, err error) {
	var pool *pgxpool.Pool
	ctx := context.Background()

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

	db := &Database{conn: pool}
	ok, err := startMigration(dsn)
	if err != nil && !ok {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping failed to database: %w", err)
	}
	repo := db.NewRepository()

	return &repo, nil
}

func startMigration(dsn string) (bool, error) {

	m, err := migrate.New("file://internal/migrations", dsn)
	if err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}
	return true, nil
}
