package bots

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type usecase interface {
	Create(ctx context.Context, name string, systemPromt string) (model.Bot, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (model.Bot, error)
	List(ctx context.Context) ([]model.Bot, error)
	Update(ctx context.Context, bot model.Bot) (model.Bot, error)
	AddProfileToBot(ctx context.Context, botID, profileID uuid.UUID) error
	RemoveProfileFromBot(ctx context.Context, botID, profileID uuid.UUID) error
	Search(ctx context.Context, query string) ([]model.Bot, error)
	GetSummary(ctx context.Context) (model.BotsSummary, error)
}

type Handler struct {
	u      usecase
	client profiles.ProfilesServiceClient
}

func NewHandler(usecase usecase, client profiles.ProfilesServiceClient) *Handler {
	return &Handler{
		u:      usecase,
		client: client,
	}
}
