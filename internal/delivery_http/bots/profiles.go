package bots

import (
	"context"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (h Handler) AddProfileToBot(ctx context.Context, params gen.AddProfileToBotParams) (gen.AddProfileToBotRes, error) {
	if err := h.u.AddProfileToBot(ctx, params.BotId, params.ProfileId); err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.AddProfileToBotNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   "bot or profile not found",
			}, nil
		}
		h.log.Error(err, true, "AddProfileToBot: add profile")
		return &gen.AddProfileToBotInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}
	return &gen.NoContent{}, nil
}

func (h Handler) RemoveProfileFromBot(ctx context.Context, params gen.RemoveProfileFromBotParams) (gen.RemoveProfileFromBotRes, error) {
	if err := h.u.RemoveProfileFromBot(ctx, params.BotId, params.ProfileId); err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.RemoveProfileFromBotNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   "bot not found",
			}, nil
		}
		h.log.Error(err, true, "RemoveProfileFromBot: remove profile")
		return &gen.RemoveProfileFromBotInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}
	return &gen.NoContent{}, nil
}

func (h Handler) GetBotProfiles(ctx context.Context, params gen.GetBotProfilesParams) (gen.GetBotProfilesRes, error) {
	_, profiles, err := h.u.GetWithProfiles(ctx, params.BotId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.GetBotProfilesNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   "bot not found",
			}, nil
		}
		h.log.Error(err, true, "GetBotProfiles: get with profiles")
		return &gen.GetBotProfilesInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}

	out := make([]gen.Profile, len(profiles))
	for i, p := range profiles {
		out[i] = gen.Profile{
			ID:           p.ID,
			Name:         p.Name,
			SystemPrompt: gen.NewOptString(p.SystemPromt),
		}
	}
	return &gen.ProfileList{Data: out}, nil
}
