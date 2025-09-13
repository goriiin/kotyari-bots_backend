package bots

import (
	"context"
	"fmt"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) Update(ctx context.Context, bot model.Bot) (model.Bot, error) {
	bot.Name = strings.TrimSpace(bot.Name)
	if bot.Name == "" {
		return model.Bot{}, fmt.Errorf("%w: name", constants.ErrValidation)
	}
	if err := s.repo.Update(ctx, bot); err != nil {
		return model.Bot{}, err
	}
	return bot, nil
}
