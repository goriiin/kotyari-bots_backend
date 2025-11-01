package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (s *Service) List(ctx context.Context) ([]model.FullBot, error) {
	bots, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	botsList := make([]model.FullBot, 0, len(bots))

	for _, bot := range bots {
		profiles, err := s.pg.GetProfilesByIDs(ctx, bot.ProfileIDs)
		if err != nil {
			return nil, err
		}
		botsList = append(botsList, model.FullBot{
			Bot:      bot,
			Profiles: profiles,
		})
	}

	return botsList, nil
}
