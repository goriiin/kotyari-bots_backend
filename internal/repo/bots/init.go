package bots

import (
	"github.com/goriiin/kotyari-bots_backend/internal/repo/pool"
)

type PGRepo struct {
	pool pool.DBPool
}

func NewPGRepo(pool pool.DBPool) *PGRepo { return &PGRepo{pool: pool} }
