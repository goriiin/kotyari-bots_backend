package bots

import (
	"context"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (h *Handler) GetBotById(ctx context.Context, params gen.GetBotByIdParams) (gen.GetBotByIdRes, error) {
	bot, profiles, err := h.u.GetWithProfiles(ctx, params.BotId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.GetBotByIdNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   err.Error(),
			}, nil
		}
		if errors.Is(err, constants.ErrServiceUnavailable) {
			return &gen.GetBotByIdInternalServerError{
				ErrorCode: constants.ServiceUnavailableMsg,
				Message:   err.Error(),
			}, nil
		}
		return nil, err
	}

	return modelToDTO(&bot, profiles), nil
}
