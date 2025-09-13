package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) List(ctx context.Context) ([]model.Bot, error) {
	return s.repo.List(ctx)
}
