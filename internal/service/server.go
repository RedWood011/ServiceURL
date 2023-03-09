package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
	"golang.org/x/exp/slog"
)

type Translation interface {
	GetURLByID(ctx context.Context, shortURL string) (string, error)
	GetAllURLsByUserID(ctx context.Context, userID string) ([]entities.URL, error)
	CreateShortURL(ctx context.Context, urls entities.URL) (ID string, err error)
	CreateShortURLs(ctx context.Context, urls []entities.URL) ([]entities.URL, error)
	PingDB(ctx context.Context) error
}

type TranslationServer struct {
	Repo    Storage
	address string
	Log     *slog.Logger
}

// New -.
func New(r Storage, log *slog.Logger, addr string) *TranslationServer {
	return &TranslationServer{
		Repo:    r,
		address: addr + "/",
		Log:     log,
	}
}
