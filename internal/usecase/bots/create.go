package bots

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"strings"
)

func (s *Service) Create(ctx context.Context, name string, systemPromt string) (model.Bot, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Bot{}, errors.Join(constants.ErrValidation, fmt.Errorf("%w: name", constants.ErrRequired))
	}
	b := model.Bot{
		ID:           uuid.New(),
		Name:         name,
		SystemPrompt: systemPromt,
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return model.Bot{}, err
	}
	return b, nil
}
