package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (h *Handler) SearchBots(ctx context.Context, params bots.SearchBotsParams) (bots.SearchBotsRes, error) {
	foundBots, err := h.u.Search(ctx, params.Q)
	if err != nil {
		h.log.Error(err, true, "SearchBots: search")
		return &bots.SearchBotsInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}

	genBots := make([]bots.Bot, len(foundBots))
	for i, b := range foundBots {
		genBots[i] = *modelToDTO(&b, nil)
	}

	return &bots.BotList{
		Data: genBots,
	}, nil
}
