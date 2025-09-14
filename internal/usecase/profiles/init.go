package profiles

import (
	"context"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type repository interface {
	Create(ctx context.Context, profile model.Profile) error
	GetByID(ctx context.Context, id uuid.UUID) (model.Profile, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error)
	List(ctx context.Context, limit int, cursor string) ([]model.Profile, error)
	Update(ctx context.Context, profile model.Profile) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{repo: r}
}
