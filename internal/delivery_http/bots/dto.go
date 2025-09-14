package bots

import (
	"github.com/google/uuid"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func modelToDTO(bot *model.Bot, profiles []gen.Profile) *gen.Bot {
	if bot == nil {
		return nil
	}
	if profiles == nil {
		profiles = []gen.Profile{}
	}

	return &gen.Bot{
		ID:                 bot.ID,
		Name:               bot.Name,
		SystemPrompt:       gen.NewOptString(bot.SystemPrompt),
		ModerationRequired: bot.ModerationRequired,
		Profiles:           profiles,
		CreatedAt:          bot.CreatedAt,
		UpdatedAt:          bot.UpdateAt,
	}
}

func dtoToModel(req *gen.BotInput, id uuid.UUID) model.Bot {
	return model.Bot{
		ID:                 id,
		Name:               req.Name,
		Email:              req.Email,
		SystemPrompt:       req.SystemPrompt.Value,
		ModerationRequired: req.ModerationRequired.Value,
		AutoPublish:        req.AutoPublish.Value,
	}
}
