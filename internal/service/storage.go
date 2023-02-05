package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

type Storage interface {
	GetFullURLByID(ctx context.Context, id string) (res string, err error)
	CreateShortURL(ctx context.Context, url []entities.URL) error
	SaveDone() error
}
