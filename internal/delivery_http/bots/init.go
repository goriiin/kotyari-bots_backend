package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type usecase interface {
	Create(ctx context.Context, bot model.Bot) (model.Bot, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithProfiles(ctx context.Context, id uuid.UUID) (model.Bot, []model.Profile, error)
	List(ctx context.Context) ([]model.FullBot, error)
	Update(ctx context.Context, bot model.Bot) (model.Bot, error)
	AddProfileToBot(ctx context.Context, botID, profileID uuid.UUID) error
	RemoveProfileFromBot(ctx context.Context, botID, profileID uuid.UUID) error
	Search(ctx context.Context, query string) ([]model.Bot, error)
	GetSummary(ctx context.Context) (model.BotsSummary, error)
}

type Handler struct {
	u   usecase
	log *logger.Logger
}

func NewHandler(usecase usecase, log *logger.Logger) *Handler {
	return &Handler{
		u:   usecase,
		log: log,
	}
}
