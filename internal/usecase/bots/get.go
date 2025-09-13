package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (model.Bot, error) {
	return s.repo.Get(ctx, id)
}
