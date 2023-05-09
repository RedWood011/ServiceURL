package postgres

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
	rows, err := r.db.Query(ctx, query, userID)
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
	query := `select short_url, original_url, user_id, is_deleted from urls where short_url = $1`
	var u entities.URL
	result := r.db.QueryRow(ctx, query, shortURL)
	if err := result.Scan(&u.ShortURL, &u.FullURL, &u.UserID, &u.IsDeleted); err != nil {
		return "", apperror.ErrDataBase
	}
	if u.IsDeleted {
		return "", apperror.ErrGone
	}
	return u.FullURL, nil
}

func (r Repository) findShortURL(ctx context.Context, fullURL string) (string, error) {
	query := `select user_id, original_url, short_url from urls where original_url = $1`
	var u entities.URL
	result := r.db.QueryRow(ctx, query, fullURL)
	if err := result.Scan(&u.UserID, &u.FullURL, &u.ShortURL); err != nil {
		return "", apperror.ErrDataBase
	}
	return u.ShortURL, nil
}

func (r Repository) CreateShortURL(ctx context.Context, url entities.URL) (string, error) {

	sqlAddRow := `INSERT INTO urls (user_id, original_url, short_url,is_deleted)
				 VALUES ($1, $2, $3, $4) `
	var pgErr *pgconn.PgError
	_, err := r.db.Exec(ctx, sqlAddRow, url.UserID, url.FullURL, url.ShortURL, url.IsDeleted)
	if err != nil {
		if errs.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			url.ShortURL, err = r.findShortURL(ctx, url.FullURL)
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperror.ErrDataBase
	}
	defer tx.Rollback(ctx)
	query := `insert into urls (short_url, original_url, user_id,is_deleted) values ($1, $2, $3, $4)`
	for _, url := range urls {
		_, err = tx.Exec(ctx, query, url.ShortURL, url.FullURL, url.UserID, url.IsDeleted)
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

func (r Repository) DeleteShortURLs(ctx context.Context, urls []string, userID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "UPDATE urls SET is_deleted = true WHERE short_url = any($1) AND user_id = $2"

	_, err = tx.Exec(ctx, query, urls, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
