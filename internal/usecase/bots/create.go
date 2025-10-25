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
		return model.Bot{}, errors.Join(constants.ErrValidation, fmt.Errorf("%w: name", constants.ErrRequired))
	}
	b := model.Bot{
		ID:                 uuid.New(),
		Name:               bot.Name,
		SystemPrompt:       bot.SystemPrompt,
		ProfileIDs:         bot.ProfileIDs,
		ModerationRequired: bot.ModerationRequired,
		CreatedAt:          time.Now(),
		UpdateAt:           time.Now(),
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return model.Bot{}, err
	}
	return b, nil
}
