package profiles

import "github.com/goriiin/kotyari-bots_backend/internal/repo/pool"

type Repository struct {
	db pool.DBPool
}

func NewRepository(db pool.DBPool) *Repository {
	return &Repository{db: db}
}
