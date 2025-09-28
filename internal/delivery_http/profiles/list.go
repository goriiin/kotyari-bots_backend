package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/profiles"
)

func (h *HTTPHandler) ListMyProfiles(ctx context.Context) (gen.ListMyProfilesRes, error) {
	profiles, err := h.u.List(ctx)
	if err != nil {
		return &gen.ListMyProfilesInternalServerError{ErrorCode: constants.ErrInternalMsg, Message: err.Error()}, nil
	}

	dtoProfiles := make([]gen.Profile, len(profiles))
	for i, p := range profiles {
		dtoProfiles[i] = *modelToHttpDTO(&p)
	}

	return &gen.ProfileList{Data: dtoProfiles}, nil
}
