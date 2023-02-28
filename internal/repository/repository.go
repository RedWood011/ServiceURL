package repository

import (
	"fmt"
	"strconv"

	"github.com/RedWood011/ServiceURL/internal/config"
	"github.com/RedWood011/ServiceURL/internal/repository/database"
	"github.com/RedWood011/ServiceURL/internal/repository/filememory"
	"github.com/RedWood011/ServiceURL/internal/repository/memory"
	"github.com/RedWood011/ServiceURL/internal/service"
)

func NewRepository(cfg *config.Config) (service.Storage, error) {
	switch {
	case cfg.FilePath != "":
		return filememory.NewFileStorage(cfg.FilePath)
	case cfg.DatabaseDSN != "":
		repetition, err := strconv.Atoi(cfg.CountRepetitionBD)
		if err != nil {
			return nil, fmt.Errorf("convert countRepetitionBD err: %w", err)
		}
		return database.NewRepo(cfg.DatabaseDSN, repetition)

	default:
		return memory.NewMemoryStorage()
	}
}
