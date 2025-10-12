package profiles

import (
	"github.com/google/uuid"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/profiles"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func modelToHttpDTO(p *model.Profile) *gen.Profile {
	if p == nil {
		return nil
	}
	return &gen.Profile{
		ID:        p.ID,
		Name:      p.Name,
		Email:     p.Email,
		Prompt:    gen.NewOptString(p.SystemPromt),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func httpDtoToModel(p *gen.ProfileInput, id uuid.UUID) model.Profile {
	return model.Profile{
		ID:          id,
		Name:        p.Name,
		Email:       p.Email,
		SystemPromt: p.Prompt.Value,
	}
}
