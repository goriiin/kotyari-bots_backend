package bots

import (
	"context"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (h *Handler) UpdateBotById(ctx context.Context, req *gen.BotInput, params gen.UpdateBotByIdParams) (gen.UpdateBotByIdRes, error) {
	_, err := h.u.Update(ctx, dtoToModel(req, params.BotId))
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.UpdateBotByIdNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   err.Error(),
			}, nil
		}
		h.log.Error(err, true, "failed to update bot")
		return &gen.UpdateBotByIdInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}

	bot, profiles, err := h.u.GetWithProfiles(ctx, params.BotId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.UpdateBotByIdNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   err.Error(),
			}, nil
		}
		h.log.Error(err, true, "failed to get updated bot")
		return &gen.UpdateBotByIdInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}

	return modelToDTO(&bot, profiles), nil
}
