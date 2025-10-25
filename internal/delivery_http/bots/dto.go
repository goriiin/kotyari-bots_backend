package bots

import (
	"github.com/google/uuid"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func modelToDTO(bot *model.Bot, profiles []model.Profile) *gen.Bot {
	if bot == nil {
		return nil
	}

	genProfiles := make([]gen.Profile, len(profiles))
	for i, p := range profiles {
		genProfiles[i] = gen.Profile{
			ID:           p.ID,
			Name:         p.Name,
			SystemPrompt: gen.NewOptString(p.SystemPromt),
		}
	}

	return &gen.Bot{
		ID:                 bot.ID,
		Name:               bot.Name,
		SystemPrompt:       gen.NewOptString(bot.SystemPrompt),
		ModerationRequired: gen.NewOptBool(bot.ModerationRequired),
		Profiles:           genProfiles,
		ProfilesCount:      len(genProfiles),
		CreatedAt:          bot.CreatedAt,
		UpdatedAt:          bot.UpdatedAt,
	}
}

func dtoToModel(req *gen.BotInput, id uuid.UUID) model.Bot {
	return model.Bot{
		ID:                 id,
		Name:               req.Name,
		SystemPrompt:       req.SystemPrompt.Value,
		ModerationRequired: req.ModerationRequired.Value,
	}
}
