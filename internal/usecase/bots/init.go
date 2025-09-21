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
}

type profileValidator interface {
	ValidateProfileExists(ctx context.Context, profileID uuid.UUID) error
}

type Service struct {
	repo repository
	pv   profileValidator
}

func NewService(r repository, pv profileValidator) *Service {
	return &Service{
		repo: r,
		pv:   pv,
	}
}
