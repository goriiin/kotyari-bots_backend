package posts

import "github.com/goriiin/kotyari-bots_backend/internal/repo/pool"

type PostsRepo struct {
	// TODO: logs
	db pool.DBPool
}

func NewAPIIntegrationsRepo(dbPool pool.DBPool) *PostsRepo {
	return &PostsRepo{
		dbPool,
	}
}
