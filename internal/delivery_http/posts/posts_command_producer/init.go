package posts_command_producer

import (
	"context"
	"time"

	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"google.golang.org/grpc"
)

type producer interface {
	Publish(ctx context.Context, env kafkaConfig.Envelope) error
	Request(ctx context.Context, env kafkaConfig.Envelope, timeout time.Duration) ([]byte, error)
}

type profilesFetcher interface {
	GetProfiles(ctx context.Context, ids []string, opts ...grpc.CallOption) (*profilesgen.GetProfilesResponse, error)
}

type botsFetcher interface {
	GetBot(ctx context.Context, id string, opts ...grpc.CallOption) (*botsgen.Bot, error)
}

type botsAndProfilesFetcher interface {
	profilesFetcher
	botsFetcher
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
