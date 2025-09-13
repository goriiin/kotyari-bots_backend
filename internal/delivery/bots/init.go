package bots

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type usecase interface {
	Create(ctx context.Context, name string, systemPromt string) (model.Bot, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (model.Bot, error)
	List(ctx context.Context) ([]model.Bot, error)
	Update(ctx context.Context, bot model.Bot) (model.Bot, error)
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

func (h *Handler) AddProfileToBot(ctx context.Context, params gen.AddProfileToBotParams) (gen.AddProfileToBotRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) CreateTaskForBotWithProfile(ctx context.Context, req *gen.TaskInput, params gen.CreateTaskForBotWithProfileParams) (gen.CreateTaskForBotWithProfileRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetBotProfiles(ctx context.Context, params gen.GetBotProfilesParams) (gen.GetBotProfilesRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetTaskById(ctx context.Context, params gen.GetTaskByIdParams) (gen.GetTaskByIdRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) RemoveProfileFromBot(ctx context.Context, params gen.RemoveProfileFromBotParams) (gen.RemoveProfileFromBotRes, error) {
	return nil, fmt.Errorf("not implemented")
}
