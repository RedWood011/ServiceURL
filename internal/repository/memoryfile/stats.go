package memoryfile

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

// GetStats
func (f *FileMap) GetStats(_ context.Context) (entities.Stats, error) {
	res := entities.Stats{
		CountUser: len(f.LongByShortURL),
		CountURL:  len(f.cacheShortURL),
	}

	return res, nil
}
