package api_integrations

import "github.com/kotyari-bots_backend/internal/repo/pool"

type APIIntegrationsRepo struct {
	db pool.DBPool
}

func NewAPIIntegrationsRepo(dbPool pool.DBPool) *APIIntegrationsRepo {
	return &APIIntegrationsRepo{
		dbPool,
	}
}
