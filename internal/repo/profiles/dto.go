package profiles

import (
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"time"
)

type profileDTO struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	SystemPromt string    `db:"system_prompt"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (d *profileDTO) toModel() model.Profile {
	return model.Profile{
		ID:          d.ID,
		Name:        d.Name,
		Email:       d.Email,
		SystemPromt: d.SystemPromt,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func toModels(dtos []profileDTO) []model.Profile {
	if dtos == nil {
		return nil
	}
	models := make([]model.Profile, len(dtos))
	for i, dto := range dtos {
		models[i] = dto.toModel()
	}
	return models
}
