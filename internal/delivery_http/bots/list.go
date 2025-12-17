package bots

import (
	"context"
	"log"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
)

func (h *Handler) ListBots(ctx context.Context) (gen.ListBotsRes, error) {
	log.Println("ListBots")
	bots, err := h.u.List(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("ListBots", bots)
	genBots := make([]gen.Bot, len(bots))
	for i, b := range bots {
		genBots[i] = *modelToDTO(&b.Bot, b.Profiles)
	}

	log.Println("bots list:", len(bots), genBots)

	return &gen.BotList{
		Data:       genBots,
		NextCursor: gen.OptNilString{},
	}, nil
}
