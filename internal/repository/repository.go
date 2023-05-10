package repository

import (
	"github.com/RedWood011/ServiceURL/internal/repository/memoryfile"
	"github.com/RedWood011/ServiceURL/internal/repository/postgres"
	"github.com/RedWood011/ServiceURL/internal/service"
)

// NewRepository Создание репозитория
func NewRepository(typeBD string, dbPostgres *postgres.Repository, dbFilebase *memoryfile.FileMap) service.Storage {
	if typeBD != "" {
		return dbPostgres
	}

	return dbFilebase
}
