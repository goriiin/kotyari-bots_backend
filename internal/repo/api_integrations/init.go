package api_integrations

import (
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/pool"
)

type APIIntegrationsRepo struct {
	log logger.Logger
	db  pool.DBPool
}

func NewAPIIntegrationsRepo(log logger.Logger, dbPool pool.DBPool) *APIIntegrationsRepo {
	return &APIIntegrationsRepo{
		log,
		dbPool,
	}
}
