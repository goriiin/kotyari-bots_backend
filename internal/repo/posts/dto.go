package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostDTO struct {
	ID        uint64      `db:"id"`
	BotID     uuid.UUID   `db:"bot_id"`
	ProfileID uuid.UUID   `db:"profile_id"`
	Platform  string      `db:"platform_type"`
	Type      pgtype.Text `db:"post_type"`
	Title     string      `db:"post_title"`
	Text      string      `db:"post_text"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

func (d PostDTO) toModel() model.Post {
	var postType model.PostType
	if d.Type.Valid {
		postType = model.PostType(d.Type.String)
	}

	return model.Post{
		ID:        d.ID,
		BotID:     d.BotID,
		ProfileID: d.ProfileID,
		Platform:  model.PlatformType(d.Platform),
		Type:      postType,
		Title:     d.Title,
		Text:      d.Text,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

type CategoryDTO struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"category_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (d CategoryDTO) toModel() model.Category {
	return model.Category{
		ID:   d.ID,
		Name: d.Name,
	}
}

func categoriesDtoToModel(categories []CategoryDTO) []model.Category {
	modelCategories := make([]model.Category, 0, len(categories))
	for _, category := range categories {
		modelCategories = append(modelCategories, category.toModel())
	}

	return modelCategories
}
