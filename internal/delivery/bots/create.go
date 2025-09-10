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

	created, err := h.u.Create(ctx, req.Name, desc)
	if err != nil {
		return nil, err
	}
	return modelToDTO(&created, nil), nil
}
