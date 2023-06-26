package postgres

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

// GetStats
func (r Repository) GetStats(ctx context.Context) (entities.Stats, error) {
	query := `SELECT COUNT(DISTINCT user_id), COUNT (DISTINCT origin_url) FROM urls;`
	row := r.DB.QueryRow(ctx, query)
	res := entities.Stats{}
	err := row.Scan(&res.CountUser, &res.CountURL)
	return res, err
}
