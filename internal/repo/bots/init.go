package bots

import (
	"github.com/goriiin/kotyari-bots_backend/internal/repo/pool"
)

type BotsRepository struct {
	db pool.DBPool
}

func NewBotsRepository(db pool.DBPool) *BotsRepository {
	return &BotsRepository{db: db}
}
