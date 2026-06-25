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
		switch {
		case errors.Is(err, constants.ErrValidation):
			return &gen.UpdateBotByIdBadRequest{
				ErrorCode: constants.ValidationMsg,
				Message:   err.Error(),
			}, nil
		case errors.Is(err, constants.ErrNotFound):
			return &gen.UpdateBotByIdNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   "bot not found",
			}, nil
		}
		h.log.Error(err, true, "UpdateBotById: update")
		return &gen.UpdateBotByIdInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   constants.InternalMsg,
		}, nil
	}

	bot, profiles, err := h.u.GetWithProfiles(ctx, params.BotId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.UpdateBotByIdNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   "bot not found",
			}, nil
		}
		h.log.Error(err, true, "UpdateBotById: get with profiles")
		return &gen.UpdateBotByIdInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   constants.InternalMsg,
		}, nil
	}

	return modelToDTO(&bot, profiles), nil
}
