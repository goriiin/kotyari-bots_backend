package profiles

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
)

func (h *GRPCHandler) GetProfiles(ctx context.Context, req *profiles.GetProfilesRequest) (*profiles.GetProfilesResponse, error) {
	profileUUIDs := make([]uuid.UUID, 0, len(req.ProfileIds))
	for _, idStr := range req.ProfileIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue // Игнорируем невалидные UUID
		}
		profileUUIDs = append(profileUUIDs, id)
	}

	if len(profileUUIDs) == 0 {
		return &profiles.GetProfilesResponse{Profiles: []*profiles.Profile{}}, nil
	}

	profileModels, err := h.u.GetByIDs(ctx, profileUUIDs)
	if err != nil {
		return nil, err
	}

	grpcProfiles := make([]*profiles.Profile, len(profileModels))
	for i, p := range profileModels {
		grpcProfiles[i] = modelToProto(&p)
	}

	return &profiles.GetProfilesResponse{Profiles: grpcProfiles}, nil
}
