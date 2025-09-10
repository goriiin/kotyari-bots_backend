package bots

import (
	"context"
	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"log"
)

func (h *Handler) UpdateBotById(ctx context.Context, req *gen.BotInput, params gen.UpdateBotByIdParams) (gen.UpdateBotByIdRes, error) {
	bot, err := h.u.Update(ctx, dtoToModel(req, params.BotId))
	if err != nil {
		return nil, err
	}

	var genProfiles []gen.Profile
	if len(bot.ProfileIDs) > 0 {
		profileIDsStr := make([]string, len(bot.ProfileIDs))
		for i, pid := range bot.ProfileIDs {
			profileIDsStr[i] = pid.String()
		}

		grpcResp, err := h.client.GetProfilesByIDs(ctx, &profiles.GetProfilesByIDsRequest{
			ProfileIds: profileIDsStr,
		})
		if err != nil {
			log.Printf("failed to get profiles for bot %s: %v", bot.ID, err)
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

	return modelToDTO(&bot, genProfiles), nil
}
