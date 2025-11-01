package bots

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (s Service) AddProfileToBot(ctx context.Context, botID, profileID uuid.UUID) error {
	// явно проверяем, что бот существует, чтобы различать 404
	if _, err := s.repo.Get(ctx, botID); err != nil {
		return err
	}
	if err := s.pv.ValidateProfilesExist(ctx, []uuid.UUID{profileID}); err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return err
		}
		return err
	}
	return s.repo.AddProfileID(ctx, botID, profileID)
}

func (s Service) RemoveProfileFromBot(ctx context.Context, botID, profileID uuid.UUID) error {
	if _, err := s.repo.Get(ctx, botID); err != nil {
		return err
	}
	return s.repo.RemoveProfileID(ctx, botID, profileID)
}
