package aggregator

import (
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/pool"
)

type AggregatorRepo struct {
	log *logger.Logger
	db  pool.DBPool
}

func NewAggregatorRepo(log *logger.Logger, dbPool pool.DBPool) *AggregatorRepo {
	return &AggregatorRepo{
		log,
		dbPool,
	}
}
