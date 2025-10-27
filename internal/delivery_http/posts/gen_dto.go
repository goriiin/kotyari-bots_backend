package posts

import (
	genCommand "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	genQuery "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

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

func QueryPostsToHttp(postsModels []model.Post) *genQuery.PostList {
	httpPosts := make([]genQuery.Post, 0, len(postsModels))
	for _, postsModel := range postsModels {
		httpPosts = append(httpPosts, *QueryModelToHttp(postsModel))
	}

	return &genQuery.PostList{
		Data: httpPosts,
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
