package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) List(ctx context.Context, limit int, cursor string) ([]model.Profile, error) {
	return s.repo.List(ctx, limit, cursor)
}
