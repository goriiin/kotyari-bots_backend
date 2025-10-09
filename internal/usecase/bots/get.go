package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) GetWithProfiles(ctx context.Context, id uuid.UUID) (model.Bot, []model.Profile, error) {
	bot, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Bot{}, nil, err
	}

	if len(bot.ProfileIDs) == 0 {
		return bot, []model.Profile{}, nil
	}

	profiles, err := s.pg.GetProfilesByIDs(ctx, bot.ProfileIDs)
	if err != nil {
		return model.Bot{}, nil, err
	}

	return bot, profiles, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (model.Bot, error) {
	return s.repo.Get(ctx, id)
}
