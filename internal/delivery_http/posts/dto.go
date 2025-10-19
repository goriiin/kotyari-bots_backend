package posts

import (
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func ModelToHttp(post model.Post) *gen.Post {
	var postType gen.OptNilPostPostType
	if post.Type != "" {
		postType = gen.NewOptNilPostPostType(gen.PostPostType(post.Type))
		postType.Null = false
	} else {
		postType.Null = true
	}

	return &gen.Post{
		ID:         post.ID,
		BotId:      post.BotID,
		ProfileId:  post.ProfileID,
		Platform:   gen.PostPlatform(post.Platform),
		PostType:   postType,
		Title:      post.Title,
		Text:       post.Text,
		Categories: nil,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}

func ModelSliceToHttpSlice(posts []model.Post) []gen.Post {
	var httpPosts []gen.Post
	for _, p := range posts {
		httpPosts = append(httpPosts, *ModelToHttp(p))
	}

	return httpPosts
}

func HttpInputToModel(input gen.PostInput) (*model.Post, string) {
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
