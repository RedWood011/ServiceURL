package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

type Storage interface {
	GetFullURLByID(ctx context.Context, shortURL string) (res string, err error)
	GetAllURLsByUserID(ctx context.Context, userID string) ([]entities.URL, error)
	CreateShortURL(ctx context.Context, url entities.URL) (string, error)
	CreateShortURLs(ctx context.Context, urls []entities.URL) ([]entities.URL, error)
	SaveDone() error
	Ping(ctx context.Context) error
}
