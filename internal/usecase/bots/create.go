package bots

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (s *Service) Create(ctx context.Context, bot model.Bot) (model.Bot, error) {
	bot.Name = strings.TrimSpace(bot.Name)
	if bot.Name == "" {
		return model.Bot{}, errors.Join(constants.ErrValidation, fmt.Errorf("name: %w", constants.ErrRequired))
	}

	if err := s.pv.ValidateProfilesExist(ctx, bot.ProfileIDs); err != nil {
		return model.Bot{}, err
	}

	now := time.Now()
	b := model.Bot{
		ID:                 uuid.New(),
		Name:               bot.Name,
		SystemPrompt:       bot.SystemPrompt,
		ModerationRequired: bot.ModerationRequired,
		ProfileIDs:         bot.ProfileIDs,
		ProfilesCount:      len(bot.ProfileIDs),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return model.Bot{}, err
	}
	return b, nil
}

// CreateWithProfiles creates a bot and returns it together with its resolved
// profiles, without re-reading the freshly created bot from the database.
// The profiles still have to be fetched because the bot only stores their IDs.
func (s *Service) CreateWithProfiles(ctx context.Context, bot model.Bot) (model.Bot, []model.Profile, error) {
	created, err := s.Create(ctx, bot)
	if err != nil {
		return model.Bot{}, nil, err
	}

	if len(created.ProfileIDs) == 0 {
		return created, []model.Profile{}, nil
	}

	profiles, err := s.pg.GetProfilesByIDs(ctx, created.ProfileIDs)
	if err != nil {
		return model.Bot{}, nil, err
	}

	return created, profiles, nil
}
