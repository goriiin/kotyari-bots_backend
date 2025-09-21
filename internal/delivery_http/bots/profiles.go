package bots

import (
	"context"
	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/pkg/ierrors"
	"log"

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
	bot, err := h.u.Get(ctx, params.BotId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.GetBotProfilesNotFound{
				ErrorCode: constants.NotFoundMsg,
				Message:   "bot not found",
			}, nil
		}
		return &gen.GetBotProfilesInternalServerError{
			ErrorCode: constants.InternalMsg,
			Message:   err.Error(),
		}, nil
	}

	if len(bot.ProfileIDs) == 0 {
		return &gen.ProfileList{
			Data: []gen.Profile{},
		}, nil
	}

	profileIDsStr := make([]string, len(bot.ProfileIDs))
	for i, pid := range bot.ProfileIDs {
		profileIDsStr[i] = pid.String()
	}

	grpcResp, err := h.client.GetProfiles(ctx, &profiles.GetProfilesRequest{
		ProfileIds: profileIDsStr,
	})
	if err != nil {
		domainErr := ierrors.GRPCToDomainError(err)
		log.Printf("failed to get profiles for bot %s: %v", bot.ID, domainErr)

		if errors.Is(domainErr, constants.ErrServiceUnavailable) {
			return &gen.GetBotProfilesInternalServerError{
				ErrorCode: constants.ServiceUnavailableMsg,
				Message:   "profiles service is unavailable",
			}, nil
		}
		return nil, domainErr
	}

	genProfiles := make([]gen.Profile, len(grpcResp.Profiles))
	for i, p := range grpcResp.Profiles {
		profileUUID, err := uuid.Parse(p.Id)
		if err != nil {
			log.Printf("error parsing profile UUID from gRPC response: %v", err)
			continue
		}
		genProfiles[i] = gen.Profile{
			ID:           profileUUID,
			Name:         p.Name,
			Email:        p.Email,
			SystemPrompt: gen.NewOptString(p.Prompt),
		}
	}

	return &gen.ProfileList{
		Data: genProfiles,
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
