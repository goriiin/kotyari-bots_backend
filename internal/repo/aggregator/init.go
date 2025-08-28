package aggregator

import "github.com/goriiin/kotyari-bots_backend/internal/repo/pool"

type AggregatorRepo struct {
	// TODO: logs
	db pool.DBPool
}

func NewAggregatorRepo(dbPool pool.DBPool) *AggregatorRepo {
	return &AggregatorRepo{
		dbPool,
	}
}
