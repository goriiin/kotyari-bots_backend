package profiles

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type usecase interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error)
	Exist(ctx context.Context, ids []uuid.UUID) (map[string]bool, error)
}

type GRPCHandler struct {
	profiles.UnimplementedProfilesServiceServer
	u   usecase
	log *logger.Logger
}

func NewGRPCHandler(u usecase, log *logger.Logger) *GRPCHandler {
	return &GRPCHandler{
		u:   u,
		log: log,
	}
}
