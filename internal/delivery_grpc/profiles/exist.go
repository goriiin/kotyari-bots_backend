package profiles

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/pkg/ierrors"
)

func (h *GRPCHandler) ProfilesExist(ctx context.Context, req *profiles.ProfilesExistRequest) (*profiles.ProfilesExistResponse, error) {
	profileUUIDs := make([]uuid.UUID, 0, len(req.ProfileIds))
	for _, idStr := range req.ProfileIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			h.log.Warn("ProfilesExist: parse id", err)
			continue
		}
		profileUUIDs = append(profileUUIDs, id)
	}

	if len(profileUUIDs) == 0 {
		return &profiles.ProfilesExistResponse{ExistenceMap: map[string]bool{}}, nil
	}

	existenceMap, err := h.u.Exist(ctx, profileUUIDs)
	if err != nil {
		h.log.Error(err, true, "ProfilesExist: exist")
		return nil, ierrors.DomainToGRPCError(err)
	}

	return &profiles.ProfilesExistResponse{ExistenceMap: existenceMap}, nil
}
