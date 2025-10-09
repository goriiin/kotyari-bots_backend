package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type repository interface {
	Create(ctx context.Context, b model.Bot) error
	Get(ctx context.Context, id uuid.UUID) (model.Bot, error)
	List(ctx context.Context) ([]model.Bot, error)
	Update(ctx context.Context, b model.Bot) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddProfileID(ctx context.Context, botID, profileID uuid.UUID) error
	RemoveProfileID(ctx context.Context, botID, profileID uuid.UUID) error
	Search(ctx context.Context, query string) ([]model.Bot, error)
	GetSummary(ctx context.Context) (model.BotsSummary, error)
}

type profileValidator interface {
	ValidateProfileExists(ctx context.Context, profileID uuid.UUID) error
}

type profileGateway interface {
	GetProfilesByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error)
}

type Service struct {
	repo repository
	pv   profileValidator
	pg   profileGateway
}

func NewService(r repository, pv profileValidator, pg profileGateway) *Service {
	return &Service{
		repo: r,
		pv:   pv,
		pg:   pg,
	}
}
