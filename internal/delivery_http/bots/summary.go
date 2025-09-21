package bots

import (
	"context"
	"github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (h *Handler) SummaryBots(ctx context.Context) (bots.SummaryBotsRes, error) {
	summary, err := h.u.GetSummary(ctx)
	if err != nil {
		return &bots.Error{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}

	return &bots.BotsSummary{
		TotalBots:             summary.TotalBots,
		TotalProfilesAttached: summary.TotalProfilesAttached,
	}, nil
}
