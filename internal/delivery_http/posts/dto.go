package posts

import (
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func modelToHttp(post model.Post) *gen.Post {
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

func modelSliceToHttpSlice(posts []model.Post) []gen.Post {
	var httpPosts []gen.Post
	for _, p := range posts {
		httpPosts = append(httpPosts, *modelToHttp(p))
	}

	return httpPosts
}

func httpInputToModel(input gen.PostInput) (*model.Post, string) {
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
