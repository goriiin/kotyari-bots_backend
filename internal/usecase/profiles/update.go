package profiles

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"strings"
	"time"
)

func (s *Service) Update(ctx context.Context, profile model.Profile) (model.Profile, error) {
	profile.Name = strings.TrimSpace(profile.Name)
	if profile.Name == "" {
		return model.Profile{}, errors.Join(constants.ErrValidation, fmt.Errorf("%w: name is required", constants.ErrRequired))
	}
	// Здесь можно добавить другие проверки (например, для email)

	existingProfile, err := s.repo.GetByID(ctx, profile.ID)
	if err != nil {
		return model.Profile{}, err
	}

	profile.CreatedAt = existingProfile.CreatedAt
	profile.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, profile); err != nil {
		return model.Profile{}, err
	}

	return profile, nil
}
