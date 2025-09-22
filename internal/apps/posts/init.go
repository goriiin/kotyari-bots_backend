package posts

import (
	"context"

	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/profiles/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_client"
	"google.golang.org/grpc"
)

type BotsProfilesFetcher interface {
	GetBot(ctx context.Context, id string, opts ...grpc.CallOption) (*botsgen.Bot, error)
	GetProfile(ctx context.Context, id string, opts ...grpc.CallOption) (*profilesgen.Profile, error)
	BatchGetProfiles(ctx context.Context, ids []string, opts ...grpc.CallOption) (*profilesgen.BatchGetProfilesResponse, error)
}

type PostsApp struct {
	fetcher BotsProfilesFetcher
	grpcCfg posts_client.PostsGRPCClientAppConfig
	appCfg  *PostsAppCfg
}

func NewPostsApp(appCfg *PostsAppCfg) (*PostsApp, error) {
	grpcClient, err := posts_client.NewPostsGRPCClient(&appCfg.GrpcClient)
	if err != nil {
		return nil, err
	}

	return &PostsApp{
		fetcher: grpcClient,
		grpcCfg: appCfg.GrpcClient,
		appCfg:  appCfg,
	}, nil
}
