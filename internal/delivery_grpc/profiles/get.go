package profiles

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type usecase interface {
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error)
}

type GRPCHandler struct {
	profiles.UnimplementedProfilesServiceServer
	u usecase
}

func NewGRPCHandler(u usecase) *GRPCHandler {
	return &GRPCHandler{u: u}
}

func (h *GRPCHandler) GetProfilesByIDs(ctx context.Context, req *profiles.GetProfilesByIDsRequest) (*profiles.GetProfilesByIDsResponse, error) {
	profileUUIDs := make([]uuid.UUID, 0, len(req.ProfileIds))
	for _, idStr := range req.ProfileIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		profileUUIDs = append(profileUUIDs, id)
	}

	profileModels, err := h.u.GetByIDs(ctx, profileUUIDs)
	if err != nil {
		return nil, err
	}

	grpcProfiles := make([]*profiles.Profile, len(profileModels))
	for i, p := range profileModels {
		grpcProfiles[i] = modelToProto(&p)
	}

	return &profiles.GetProfilesByIDsResponse{Profiles: grpcProfiles}, nil
}
