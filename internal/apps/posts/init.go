package posts

import (
	"context"
	"fmt"

	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/profiles/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/grok_client"
	"google.golang.org/grpc"
)

type BotsProfilesFetcher interface {
	GetBot(ctx context.Context, id string, opts ...grpc.CallOption) (*botsgen.Bot, error)
	GetProfile(ctx context.Context, id string, opts ...grpc.CallOption) (*profilesgen.Profile, error)
	BatchGetProfiles(ctx context.Context, ids []string, opts ...grpc.CallOption) (*profilesgen.BatchGetProfilesResponse, error)
}

type PostGenerator interface {
	GeneratePost(ctx context.Context, botPrompt, profilePrompt string) (string, error)
}

type PostsApp struct {
	fetcher   BotsProfilesFetcher
	appCfg    *PostsAppCfg
	generator PostGenerator
}

func NewPostsApp(appCfg *PostsAppCfg) (*PostsApp, error) {
	grpcClient, err := posts_client.NewPostsGRPCClient(&appCfg.GrpcClient)
	if err != nil {
		return nil, err
	}

	grokClient, err := grok_client.NewGrokClient(&appCfg.GrokCfg)
	if err != nil {
		return nil, err
	}

	// Тестовый запрос, будет убран в будущем
	post, err := grokClient.GeneratePost(context.Background(), "You are a test assistant.", "Testing. Just say hi and hello world and nothing else.")
	if err != nil {
		return nil, err
	}
	fmt.Println(post)

	return &PostsApp{
		fetcher:   grpcClient,
		appCfg:    appCfg,
		generator: grokClient,
	}, nil
}
