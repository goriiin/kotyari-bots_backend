package posts_command_producer

import (
	"context"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
)

func (p *PostsCommandHandler) CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error) {
	//bot, err := p.fetcher.GetBot(ctx, req.BotId.String())
	//if err != nil {
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}
	//mockedBot := struct {
	//	Id        string
	//	BotPrompt string
	//}{
	//	uuid.New().String(),
	//	"You are a test assistant.",
	//}
	//
	//mockedProfiles := []struct {
	//	Id            string
	//	ProfilePrompt string
	//}{
	//	{
	//		uuid.New().String(),
	//		"Testing. Just say hi and hello world and nothing else.",
	//	},
	//	{
	//		uuid.New().String(),
	//		"Add a word test in the end",
	//	},
	//}

	//idsString := make([]string, 0, len(req.ProfileIds))
	//for _, id := range req.ProfileIds {
	//	idsString = append(idsString, id.String())
	//}
	//
	//profilesBatch, err := p.fetcher.BatchGetProfiles(ctx, idsString)
	//if err != nil {
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}

	// -> kafka

	return &gen.PostList{
		Data: nil,
	}, nil
}
