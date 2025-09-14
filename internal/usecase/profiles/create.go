package profiles

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"strings"
)

func (s *Service) Create(ctx context.Context, profile model.Profile) (model.Profile, error) {
	profile.Name = strings.TrimSpace(profile.Name)
	if profile.Name == "" {
		return model.Profile{}, errors.Join(constants.ErrValidation, fmt.Errorf("%w: name", constants.ErrRequired))
	}
	profile.ID = uuid.New()

	if err := s.repo.Create(ctx, profile); err != nil {
		return model.Profile{}, err
	}
	return profile, nil
}
