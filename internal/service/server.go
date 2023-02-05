package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

type Translation interface {
	GetURLByID(ctx context.Context, id string) (string, error)
	CreateShortURL(ctx context.Context, urls []entities.URL) (IDs []string, err error)
}

type TranslationServer struct {
	Repo    Storage
	address string
}

// New -.
func New(r Storage, addr string) *TranslationServer {
	return &TranslationServer{
		Repo:    r,
		address: addr,
	}
}
