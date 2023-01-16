package repository

import (
	"github.com/RedWood011/ServiceURL/internal/repository/cache"
	"github.com/RedWood011/ServiceURL/internal/service"
)

func NewRepository(db string) service.Storage {
	switch db {
	default:
		return cache.NewURLStorage()
	}
}
