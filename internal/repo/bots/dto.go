package bots

import (
	"time"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type botDTO struct {
	ID                 uuid.UUID   `db:"id"`
	Name               string      `db:"bot_name"`
	SystemPrompt       string      `db:"system_prompt"`
	ModerationRequired bool        `db:"moderation_required"`
	ProfileIDs         []uuid.UUID `db:"profile_ids"`
	ProfilesCount      int         `db:"profiles_count"`
	CreatedAt          time.Time   `db:"created_at"`
	UpdatedAt          time.Time   `db:"updated_at"`
}

func (d botDTO) toModel() model.Bot {
	return model.Bot{
		ID:                 d.ID,
		Name:               d.Name,
		SystemPrompt:       d.SystemPrompt,
		ModerationRequired: d.ModerationRequired,
		ProfileIDs:         d.ProfileIDs,
		ProfilesCount:      d.ProfilesCount,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}

func toModels(dtos []botDTO) []model.Bot {
	if dtos == nil {
		return nil
	}
	models := make([]model.Bot, len(dtos))
	for i, dto := range dtos {
		models[i] = dto.toModel()
	}
	return models
}
