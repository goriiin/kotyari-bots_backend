package bots

import (
	"context"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
)

func (h *Handler) ListBots(ctx context.Context) (gen.ListBotsRes, error) {
	bots, err := h.u.List(ctx)
	if err != nil {
		return nil, err
	}

	genBots := make([]gen.Bot, len(bots))
	for i, b := range bots {
		genBots[i] = *modelToDTO(&b, nil)
	}

	return &gen.BotList{
		Data:       genBots,
		NextCursor: gen.OptNilString{},
	}, nil
}
