package bots

import (
	"context"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) GetSummary(ctx context.Context) (model.BotsSummary, error) {
	return s.repo.GetSummary(ctx)
}
