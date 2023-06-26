package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

// GetAllStats
func (s *TranslationServer) GetAllStats(ctx context.Context) (entities.Stats, error) {
	return s.Repo.GetStats(ctx)
}
