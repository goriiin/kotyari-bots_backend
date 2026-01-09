package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostCheckDTO struct {
	ID      uuid.UUID `db:"id"`
	GroupID uuid.UUID `db:"group_id"`
	Title   string    `db:"post_title"`
	Text    string    `db:"post_text"`
	IsSeen  bool      `db:"is_seen"`
}

func (c PostCheckDTO) ToModel() model.Post {
	return model.Post{
		ID:      c.ID,
		GroupID: c.GroupID,
		Title:   c.Title,
		Text:    c.Text,
		IsSeen:  c.IsSeen,
	}
}

func PostCheckDTOToModelSlice(dto []PostCheckDTO) []model.Post {
	posts := make([]model.Post, 0, len(dto))
	for _, postCheckDTO := range dto {
		posts = append(posts, postCheckDTO.ToModel())
	}

	return posts
}

type PostDTO struct {
	ID          uuid.UUID     `db:"id"`
	OtvetiID    pgtype.Uint64 `db:"otveti_id"`
	BotID       uuid.UUID     `db:"bot_id"`
	BotName     string        `db:"bot_name"`
	ProfileID   uuid.UUID     `db:"profile_id"`
	ProfileName string        `db:"profile_name"`
	GroupID     uuid.UUID     `db:"group_id"`
	UserPrompt  string        `db:"user_prompt"`
	Platform    string        `db:"platform_type"`
	Type        pgtype.Text   `db:"post_type"`
	Title       string        `db:"post_title"`
	Text        string        `db:"post_text"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}

func (d PostDTO) ToModel() model.Post {
	var postType model.PostType
	if d.Type.Valid {
		postType = model.PostType(d.Type.String)
	}

	return model.Post{
		ID:          d.ID,
		GroupID:     d.GroupID,
		OtvetiID:    d.OtvetiID.Uint64,
		BotID:       d.BotID,
		BotName:     d.BotName,
		ProfileID:   d.ProfileID,
		ProfileName: d.ProfileName,
		Platform:    model.PlatformType(d.Platform),
		Type:        postType,
		UserPrompt:  d.UserPrompt,
		Title:       d.Title,
		Text:        d.Text,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func PostsDTOToModel(postsDTO []PostDTO) []model.Post {
	postsModel := make([]model.Post, 0, len(postsDTO))
	for _, postDTO := range postsDTO {
		postsModel = append(postsModel, postDTO.ToModel())
	}

	return postsModel
}

type CategoryDTO struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"category_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (d CategoryDTO) ToModel() model.Category {
	return model.Category{
		ID:   d.ID,
		Name: d.Name,
	}
}

func CategoriesDtoToModel(categories []CategoryDTO) []model.Category {
	modelCategories := make([]model.Category, 0, len(categories))
	for _, category := range categories {
		modelCategories = append(modelCategories, category.ToModel())
	}

	return modelCategories
}
