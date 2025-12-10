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
		ID:          post.ID,
		GroupId:     post.GroupID,
		BotId:       post.BotID,
		BotName:     post.BotName,
		ProfileId:   post.ProfileID,
		ProfileName: post.ProfileName,
		Platform:    genQuery.PostPlatform(post.Platform),
		PostType:    postType,
		Title:       post.Title,
		Text:        post.Text,
		Categories:  nil,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
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
		ID:          post.ID,
		GroupId:     post.GroupID,
		BotId:       post.BotID,
		BotName:     post.BotName,
		ProfileId:   post.ProfileID,
		ProfileName: post.ProfileName,
		Platform:    genCommand.PostPlatform(post.Platform),
		PostType:    postType,
		Title:       post.Title,
		Text:        post.Text,
		Categories:  nil,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}

func PostsToCandidates(posts []model.Post) []model.Candidate {
	candidates := make([]model.Candidate, 0, len(posts))
	for _, post := range posts {
		candidates = append(candidates, model.Candidate{
			Title: post.Title,
			Text:  post.Text,
		})
	}

	return candidates
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

func PostsCheckModelToHttp(post model.Post) genQuery.PostsCheckObject {
	isReady := !(post.Text == "" || post.Title == "")

	return genQuery.PostsCheckObject{
		ID:      post.ID,
		GroupID: post.GroupID,
		IsReady: isReady,
	}
}

func PostsCheckModelsToHttpSlice(posts []model.Post) *genQuery.PostsCheckList {
	checkObjects := make([]genQuery.PostsCheckObject, 0, len(posts))

	for _, post := range posts {
		checkObjects = append(checkObjects, PostsCheckModelToHttp(post))
	}

	return &genQuery.PostsCheckList{
		Data: checkObjects,
	}
}
