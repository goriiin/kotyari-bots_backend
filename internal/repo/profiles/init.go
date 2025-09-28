package profiles

import "github.com/goriiin/kotyari-bots_backend/internal/repo/pool"

type Repository struct {
	pool pool.DBPool
}

func NewRepository(pool pool.DBPool) *Repository {
	return &Repository{pool: pool}
}
