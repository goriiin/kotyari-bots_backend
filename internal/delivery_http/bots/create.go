package bots

import (
	"context"

	"github.com/google/uuid"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (h *Handler) CreateBot(ctx context.Context, req *gen.BotInput) (gen.CreateBotRes, error) {
	var desc string
	if v, ok := req.SystemPrompt.Get(); ok {
		desc = v
	}

	var moderation bool
	if v, ok := req.ModerationRequired.Get(); ok {
		moderation = v
	}

	profiles := make([]uuid.UUID, 0, len(req.Profiles))
	for _, p := range req.Profiles {
		profiles = append(profiles, p.ID)
	}

	created, err := h.u.Create(ctx, model.Bot{
		Name:               req.Name,
		SystemPrompt:       desc,
		ModerationRequired: moderation,
		ProfileIDs:         profiles,
	})
	if err != nil {
		return nil, err
	}
	return modelToDTO(&created, nil), nil
}
