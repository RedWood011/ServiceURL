package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

type Translation interface {
	GetURLByID(ctx context.Context, id string) (string, error)
	CreateShortURL(ctx context.Context, urls []entities.URL) (IDs []string, err error)
}
type Storage interface {
	GetFullURLByID(ctx context.Context, id string) (res string, err error)
	GetIDsByURLs(ctx context.Context, urls []string) (map[string]string, error)
	CreateShortURL(ctx context.Context, url []entities.URL) error
}

// TranslationUseCase -.
type TranslationServer struct {
	repo    Storage
	address string
}

// New -.
func New(r Storage, addr string) *TranslationServer {
	return &TranslationServer{
		repo:    r,
		address: addr,
	}
}
