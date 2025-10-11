package bots

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (s *Service) AddProfileToBot(ctx context.Context, botID, profileID uuid.UUID) error {
	if _, err := s.repo.Get(ctx, botID); err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return errors.Wrapf(err, "bot with id %s not found", botID)
		}
		return err
	}

	if err := s.pv.ValidateProfileExists(ctx, profileID); err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return errors.Wrapf(err, "profile with id %s not found or not accessible", profileID)
		}
		return err
	}

	return s.repo.AddProfileID(ctx, botID, profileID)
}

func (s *Service) RemoveProfileFromBot(ctx context.Context, botID, profileID uuid.UUID) error {
	if _, err := s.repo.Get(ctx, botID); err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return errors.Wrapf(err, "bot with id %s not found", botID)
		}
		return err
	}

	return s.repo.RemoveProfileID(ctx, botID, profileID)
}
