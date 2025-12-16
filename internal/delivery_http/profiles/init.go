package profiles

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type usecase interface {
	Create(ctx context.Context, profile model.Profile) (model.Profile, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Profile, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error)
	List(ctx context.Context) ([]model.Profile, error)
	Update(ctx context.Context, profile model.Profile) (model.Profile, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type HTTPHandler struct {
	u   usecase
	log *logger.Logger
}

func NewHTTPHandler(u usecase, log *logger.Logger) *HTTPHandler {
	return &HTTPHandler{
		u:   u,
		log: log,
	}
}
