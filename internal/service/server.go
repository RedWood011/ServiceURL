package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/RedWood011/ServiceURL/internal/workers"
	"golang.org/x/exp/slog"
)

// Translation Интерфейс взаимодействия с Repository
type Translation interface {
	GetURLByID(ctx context.Context, shortURL string) (string, error)
	GetAllURLsByUserID(ctx context.Context, userID string) ([]entities.URL, error)
	CreateShortURL(ctx context.Context, urls entities.URL) (ID string, err error)
	CreateShortURLs(ctx context.Context, urls []entities.URL) ([]entities.URL, error)
	PingDB(ctx context.Context) error
	DeleteShortURLs(ctx context.Context, urls []string, usedID string)
	GetAllStats(ctx context.Context) (entities.Stats, error)
}

// TranslationServer Сервер
type TranslationServer struct {
	Repo    Storage
	address string
	Log     *slog.Logger
	wp      *workers.WorkerPool
}

// New -.
func New(r Storage, log *slog.Logger, wp *workers.WorkerPool, addr string) *TranslationServer {
	return &TranslationServer{
		Repo:    r,
		address: addr + "/",
		Log:     log,
		wp:      wp,
	}
}
