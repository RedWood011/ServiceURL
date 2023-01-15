package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

type Translation interface {
	GetUrlByID(ctx context.Context, id string) (string, error)
	CreateShortUrl(ctx context.Context, urls []entities.Url) (IDs []string, err error)
}
type Storage interface {
	GetFullUrlByID(ctx context.Context, id string) (res string, err error)
	GetIDsByUrls(ctx context.Context, urls []string) (map[string]string, error)
	CreateShortUrl(ctx context.Context, url []entities.Url) error
}

// TranslationUseCase -.
type TranslationServer struct {
	repo Storage
}

// New -.
func New(r Storage) *TranslationServer {
	return &TranslationServer{
		repo: r,
	}
}
