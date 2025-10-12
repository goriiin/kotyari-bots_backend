package profiles

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type Usecase interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error)
	Exist(ctx context.Context, ids []uuid.UUID) (map[string]bool, error)
}

type GRPCHandler struct {
	profiles.UnimplementedProfilesServiceServer
	u Usecase
}

func NewGRPCHandler(u Usecase) *GRPCHandler {
	return &GRPCHandler{u: u}
}
