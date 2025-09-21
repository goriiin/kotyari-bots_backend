package bots

import (
	"context"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (h *Handler) AddProfileToBot(ctx context.Context, params gen.AddProfileToBotParams) (gen.AddProfileToBotRes, error) {
	err := h.u.AddProfileToBot(ctx, params.BotId, params.ProfileId)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrNotFound):
			return &gen.AddProfileToBotNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   err.Error(),
			}, nil
		case errors.Is(err, constants.ErrValidation):
			return &gen.AddProfileToBotBadRequest{
				ErrorCode: constants.ValidationMsg,
				Message:   err.Error(),
			}, nil
		case errors.Is(err, constants.ErrServiceUnavailable):
			return &gen.AddProfileToBotInternalServerError{
				ErrorCode: constants.ServiceUnavailableMsg,
				Message:   err.Error(),
			}, nil
		default:
			return nil, err
		}
	}
	return &gen.NoContent{}, nil
}

func (h *Handler) GetBotProfiles(ctx context.Context, params gen.GetBotProfilesParams) (gen.GetBotProfilesRes, error) {
	return &gen.GetBotProfilesNotFound{
		ErrorCode: constants.NotImplementedMsg,
		Message:   "not implemented",
	}, nil
}

func (h *Handler) RemoveProfileFromBot(ctx context.Context, params gen.RemoveProfileFromBotParams) (gen.RemoveProfileFromBotRes, error) {
	err := h.u.RemoveProfileFromBot(ctx, params.BotId, params.ProfileId)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrNotFound):
			return &gen.RemoveProfileFromBotNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   err.Error(),
			}, nil
		case errors.Is(err, constants.ErrServiceUnavailable):
			return &gen.RemoveProfileFromBotInternalServerError{
				ErrorCode: constants.ServiceUnavailableMsg,
				Message:   err.Error(),
			}, nil
		default:
			return nil, err
		}
	}
	return &gen.NoContent{}, nil
}
