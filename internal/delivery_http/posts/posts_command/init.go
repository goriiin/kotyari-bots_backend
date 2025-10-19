package posts_command

import (
	"context"

	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	postssgen "github.com/goriiin/kotyari-bots_backend/api/protos/posts/gen"
	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/profiles/gen"
	"google.golang.org/grpc"
)

type postsGenerator interface {
	GetPost(ctx context.Context, userPrompt, profilePrompt, botPrompt string, opts ...grpc.CallOption) (*postssgen.GetPostResponse, error)
}

type profilesFetcher interface {
	GetProfile(ctx context.Context, id string, opts ...grpc.CallOption) (*profilesgen.Profile, error)
	BatchGetProfiles(ctx context.Context, ids []string, opts ...grpc.CallOption) (*profilesgen.BatchGetProfilesResponse, error)
}

type botsFetcher interface {
	GetBot(ctx context.Context, id string, opts ...grpc.CallOption) (*botsgen.Bot, error)
}

type botsAndProfilesFetcher interface {
	profilesFetcher
	botsFetcher
}

type PostsCommandHandler struct {
	generator postsGenerator
	fetcher   botsAndProfilesFetcher
}

func NewPostsHandler(generator postsGenerator, fetcher botsAndProfilesFetcher) *PostsCommandHandler {
	return &PostsCommandHandler{
		generator: generator,
		fetcher:   fetcher,
	}
}
