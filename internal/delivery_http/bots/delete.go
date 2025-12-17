package bots

import (
	"context"
	"log"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
)

func (h *Handler) DeleteBotById(ctx context.Context, params gen.DeleteBotByIdParams) (gen.DeleteBotByIdRes, error) {
	log.Println("delete:", params.BotId)
	err := h.u.Delete(ctx, params.BotId)
	if err != nil {
		return nil, err
	}

	return &gen.NoContent{}, nil
}
