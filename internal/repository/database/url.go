package database

import (
	"context"

	errs "errors"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func (r Repository) GetAllURLsByUserID(ctx context.Context, userID string) ([]entities.URL, error) {
	query := `select short_url, original_url, user_id from urls where user_id = $1`
	var result []entities.URL
	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, apperror.ErrDataBase
	}

	if rows.Err() != nil {
		return result, rows.Err()
	}

	defer rows.Close()

	for rows.Next() {
		var u entities.URL
		err = rows.Scan(&u.ShortURL, &u.FullURL, &u.UserID)
		if err != nil {
			return nil, apperror.ErrDataBase
		}
		result = append(result, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.ErrDataBase
	}

	return result, nil
}

func (r Repository) GetFullURLByID(ctx context.Context, shortURL string) (res string, err error) {
	query := `select short_url, original_url, user_id from urls where short_url = $1`
	var u entities.URL
	result := r.conn.QueryRow(ctx, query, shortURL)
	if err := result.Scan(&u.ShortURL, &u.FullURL, &u.UserID); err != nil {
		return "", apperror.ErrDataBase
	}
	return u.FullURL, nil
}

func (r Repository) findShortUrl(ctx context.Context, fullURL string) (string, error) {
	query := `select user_id, original_url, short_url from urls where original_url = $1`
	var u entities.URL
	result := r.conn.QueryRow(ctx, query, fullURL)
	if err := result.Scan(&u.UserID, &u.FullURL, &u.ShortURL); err != nil {
		return "", apperror.ErrDataBase
	}
	return u.ShortURL, nil
}

func (r Repository) CreateShortURL(ctx context.Context, url entities.URL) (string, error) {

	sqlAddRow := `INSERT INTO urls (user_id, original_url, short_url)
				 VALUES ($1, $2, $3) `
	var pgErr *pgconn.PgError
	_, err := r.conn.Exec(ctx, sqlAddRow, url.UserID, url.FullURL, url.ShortURL)
	if err != nil {
		if errs.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			url.ShortURL, err = r.findShortUrl(ctx, url.FullURL)
			if err != nil {
				return "", apperror.ErrDataBase
			}
			return url.ShortURL, apperror.ErrConflict

		} else {
			return "", apperror.ErrDataBase
		}
	}

	return url.ShortURL, nil
}

func (r Repository) CreateShortURLs(ctx context.Context, urls []entities.URL) ([]entities.URL, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, apperror.ErrDataBase
	}
	defer tx.Rollback(ctx)
	query := `insert into urls (short_url, original_url, user_id) values ($1, $2, $3)`
	for _, url := range urls {
		_, err = tx.Exec(context.Background(), query, url.ShortURL, url.FullURL, url.UserID)
		if err != nil {
			return nil, apperror.ErrDataBase
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, apperror.ErrDataBase
	}
	return urls, nil
}

func (r Repository) SaveDone() error {
	r.conn.Close()
	return nil
}

func (r Repository) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}
