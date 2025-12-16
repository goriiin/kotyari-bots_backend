package bots

import (
	"context"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (h *Handler) DeleteBotById(ctx context.Context, params gen.DeleteBotByIdParams) (gen.DeleteBotByIdRes, error) {
	err := h.u.Delete(ctx, params.BotId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.DeleteBotByIdNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   err.Error(),
			}, nil
		}
		h.log.Error(err, true, "failed to delete bot")
		return &gen.DeleteBotByIdInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}
	return &gen.NoContent{}, nil
}
