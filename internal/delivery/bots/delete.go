package bots

import (
	"context"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
)

func (h *Handler) DeleteBotById(ctx context.Context, params gen.DeleteBotByIdParams) (gen.DeleteBotByIdRes, error) {
	err := h.u.Delete(ctx, params.BotId)
	if err != nil {
		return nil, err
	}
	return &gen.NoContent{}, nil
}
