package bots

import (
	"context"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
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

	created, err := h.u.Create(ctx, req.Name, desc, moderation)
	if err != nil {
		return nil, err
	}
	return modelToDTO(&created, nil), nil
}
