package bots

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
)

func (h *Handler) GetBotById(ctx context.Context, params gen.GetBotByIdParams) (gen.GetBotByIdRes, error) {
	b, err := h.u.Get(ctx, params.BotId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.GetBotByIdNotFound{
				ErrorCode: constants.ErrNotFoundMsg,
				Message:   "bot not found",
			}, nil
		}
		return nil, err
	}

	var genProfiles []gen.Profile
	if len(b.ProfileIDs) > 0 {
		profileIDsStr := make([]string, len(b.ProfileIDs))
		for i, pid := range b.ProfileIDs {
			profileIDsStr[i] = pid.String()
		}

		grpcResp, err := h.client.GetProfilesByIDs(ctx, &profiles.GetProfilesByIDsRequest{
			ProfileIds: profileIDsStr,
		})
		if err != nil {
			return nil, err
		} else {
			genProfiles = make([]gen.Profile, len(grpcResp.Profiles))
			for i, p := range grpcResp.Profiles {
				profileUUID, _ := uuid.Parse(p.Id)
				genProfiles[i] = gen.Profile{
					ID:           profileUUID,
					Name:         p.Name,
					Email:        p.Email,
					SystemPrompt: gen.NewOptString(p.Prompt),
				}
			}
		}
	}

	return modelToDTO(&b, genProfiles), nil
}
