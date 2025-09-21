package bots

import (
	"context"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) Search(ctx context.Context, query string) ([]model.Bot, error) {
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		return []model.Bot{}, nil
	}
	return s.repo.Search(ctx, trimmedQuery)
}
