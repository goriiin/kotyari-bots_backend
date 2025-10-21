package posts_command_producer

import (
	"context"
	"time"

	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	postssgen "github.com/goriiin/kotyari-bots_backend/api/protos/posts/gen"
	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/profiles/gen"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"google.golang.org/grpc"
)

type producer interface {
	Publish(ctx context.Context, env kafkaConfig.Envelope) error
	Request(ctx context.Context, env kafkaConfig.Envelope, timeout time.Duration) ([]byte, error)
}

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
	postsGenerator
}

type PostsCommandHandler struct {
	fetcher  botsAndProfilesFetcher
	producer producer
}

func NewPostsHandler(fetcher botsAndProfilesFetcher, producer producer) *PostsCommandHandler {
	return &PostsCommandHandler{
		fetcher:  fetcher,
		producer: producer,
	}
}
