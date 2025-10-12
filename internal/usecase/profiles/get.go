package profiles

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error) {
	return s.repo.GetByIDs(ctx, ids)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (model.Profile, error) {
	return s.repo.GetByID(ctx, id)
}
