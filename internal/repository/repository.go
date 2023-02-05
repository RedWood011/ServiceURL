package repository

import (
	"github.com/RedWood011/ServiceURL/internal/config"
	"github.com/RedWood011/ServiceURL/internal/repository/cache"
	"github.com/RedWood011/ServiceURL/internal/repository/filememory"
	"github.com/RedWood011/ServiceURL/internal/service"
)

func NewRepository(cfg *config.Config) (service.Storage, error) {
	switch {
	case cfg.FilePath != "":
		return filememory.NewFileStorage(cfg.FilePath)

	default:
		return cache.NewMemoryStorage()
	}

}
