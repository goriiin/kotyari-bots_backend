package posts

import (
	genCommand "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	genQuery "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

const (
	CmdCreate kafkaConfig.Command = "create"
	CmdUpdate kafkaConfig.Command = "update"
	CmdDelete kafkaConfig.Command = "delete"
)

type KafkaResponse struct {
	Status string `json:"status"`
	// RAW пост, пока пусть будет OK
}

func PayloadToEnvelope(command kafkaConfig.Command, entityID string, payload []byte) kafkaConfig.Envelope {
	return kafkaConfig.Envelope{
		Command:       command,
		EntityID:      entityID,
		Payload:       nil,
		CorrelationID: "",
		Attempt:       0,
	}
}

type RawPostUpdate struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func ToRawPost(post model.Post) RawPostUpdate {
	return RawPostUpdate{
		ID:    post.ID,
		Title: post.Title,
		Text:  post.Text,
	}
}

func FromRawPost(rawPost RawPostUpdate) model.Post {
	return model.Post{
		ID:    rawPost.ID,
		Title: rawPost.Title,
		Text:  rawPost.Title,
	}
}

func QueryModelToHttp(post model.Post) *genQuery.Post {
	var postType genQuery.OptNilPostPostType
	if post.Type != "" {
		postType = genQuery.NewOptNilPostPostType(genQuery.PostPostType(post.Type))
		postType.Null = false
	} else {
		postType.Null = true
	}

	return &genQuery.Post{
		ID:         post.ID,
		BotId:      post.BotID,
		ProfileId:  post.ProfileID,
		Platform:   genQuery.PostPlatform(post.Platform),
		PostType:   postType,
		Title:      post.Title,
		Text:       post.Text,
		Categories: nil,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}

func ModelToHttp(post model.Post) *genCommand.Post {
	var postType genCommand.OptNilPostPostType
	if post.Type != "" {
		postType = genCommand.NewOptNilPostPostType(genCommand.PostPostType(post.Type))
		postType.Null = false
	} else {
		postType.Null = true
	}

	return &genCommand.Post{
		ID:         post.ID,
		BotId:      post.BotID,
		ProfileId:  post.ProfileID,
		Platform:   genCommand.PostPlatform(post.Platform),
		PostType:   postType,
		Title:      post.Title,
		Text:       post.Text,
		Categories: nil,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}

func ModelSliceToHttpSlice(posts []model.Post) []genCommand.Post {
	var httpPosts []genCommand.Post
	for _, p := range posts {
		httpPosts = append(httpPosts, *ModelToHttp(p))
	}

	return httpPosts
}

func HttpInputToModel(input genCommand.PostInput) (*model.Post, string) {
	var postType model.PostType
	if !input.PostType.Null {
		postType = model.PostType(input.PostType.Value)
	}

	return &model.Post{
		BotID:    input.BotId,
		Platform: model.PlatformType(input.Platform),
		Type:     postType,
	}, input.TaskText
}
