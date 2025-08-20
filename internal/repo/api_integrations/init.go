package api_integrations

import "github.com/goriiin/kotyari-bots_backend/internal/repo/pool"

type APIIntegrationsRepo struct {
	// TODO: logs
	db pool.DBPool
}

func NewAPIIntegrationsRepo(dbPool pool.DBPool) *APIIntegrationsRepo {
	return &APIIntegrationsRepo{
		dbPool,
	}
}
