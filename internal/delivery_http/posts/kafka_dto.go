package posts

import (
	"github.com/google/uuid"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

const (
	CmdCreate kafkaConfig.Command = "create"
	CmdUpdate kafkaConfig.Command = "update"
	CmdDelete kafkaConfig.Command = "delete"
)

// KafkaResponse TODO: model.Post -> []model.Post?
type KafkaResponse struct {
	Error string     `json:"error,omitempty"`
	Post  model.Post `json:"post"`
}

type KafkaDeletePostRequest struct {
	PostID uuid.UUID `json:"post_id"`
}

type KafkaCreatePostRequest struct {
	PostID     uuid.UUID            `json:"post_id"`
	BotID      uuid.UUID            `json:"bot_id"`
	UserPrompt string               `json:"user_prompt"`
	BotPrompt  string               `json:"bot_prompt"`
	Profiles   []CreatePostProfiles `json:"profiles"`
	Platform   model.PlatformType   `json:"platform_type"`
	PostType   model.PostType       `json:"post_type"`
}

type CreatePostProfiles struct {
	ProfileID     uuid.UUID `json:"profile_id"`
	ProfilePrompt string    `json:"profile_prompt"`
}

type KafkaUpdatePostRequest struct {
	PostID uuid.UUID `json:"post_id"`
	Title  string    `json:"title"`
	Text   string    `json:"text"`
}

func PayloadToEnvelope(command kafkaConfig.Command, entityID string, payload []byte) kafkaConfig.Envelope {
	return kafkaConfig.Envelope{
		Command:  command,
		EntityID: entityID,
		Payload:  payload,
	}
}

func (r KafkaResponse) PostCommandToGen() *gen.Post {
	postType := gen.OptNilPostPostType{
		Value: gen.PostPostType(r.Post.Type),
		Set:   true,
	}

	return &gen.Post{
		ID:         r.Post.ID,
		OtvetiId:   r.Post.OtvetiID,
		BotId:      r.Post.BotID,
		ProfileId:  r.Post.ProfileID,
		Platform:   gen.PostPlatform(r.Post.Platform),
		PostType:   postType,
		Title:      r.Post.Title,
		Text:       r.Post.Text,
		Categories: nil, // TODO: ??
		CreatedAt:  r.Post.CreatedAt,
		UpdatedAt:  r.Post.UpdatedAt,
	}
}

func (u KafkaUpdatePostRequest) ToModel() model.Post {
	return model.Post{
		ID:    u.PostID,
		Title: u.Title,
		Text:  u.Text,
	}
}
