package bots

import (
	"context"
	"fmt"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (s *Service) Update(ctx context.Context, bot model.Bot) (model.Bot, error) {
	bot.Name = strings.TrimSpace(bot.Name)
	if bot.Name == "" {
		return model.Bot{}, fmt.Errorf("name: %w", constants.ErrValidation)
	}
	if err := s.pv.ValidateProfilesExist(ctx, bot.ProfileIDs); err != nil {
		return model.Bot{}, err
	}
	bot.ProfilesCount = len(bot.ProfileIDs)
	if err := s.repo.Update(ctx, bot); err != nil {
		return model.Bot{}, err
	}
	return bot, nil
}
