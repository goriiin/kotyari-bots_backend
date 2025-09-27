package profiles

import (
	"context"
	"github.com/google/uuid"
)

// Реализация бизнес-логики для проверки существования профилей
func (s *Service) Exist(ctx context.Context, ids []uuid.UUID) (map[string]bool, error) {
	return s.repo.Exist(ctx, ids)
}
