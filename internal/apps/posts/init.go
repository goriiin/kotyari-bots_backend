package posts

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/grok_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_command"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts"
	postsRepoLib "github.com/goriiin/kotyari-bots_backend/internal/repo/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"github.com/goriiin/kotyari-bots_backend/pkg/proxy"
)

type postsHttpHandler interface {
	CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error)
	CreatePostSEO(ctx context.Context, req *gen.PostInput) (gen.CreatePostSEORes, error)
	GetPostById(ctx context.Context, params gen.GetPostByIdParams) (gen.GetPostByIdRes, error)
	UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error)
	DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error)
	ListPosts(ctx context.Context) (gen.ListPostsRes, error)
}

type PostsApp struct {
	appCfg *PostsAppCfg
	grpc   *posts_client.PostsGRPCClient
	http   postsHttpHandler
}

func NewPostsApp(appCfg *PostsAppCfg, proxyCfg *proxy.ProxyConfig) (*PostsApp, error) {
	//grpcClient, err := posts_client.NewPostsGRPCClient(&appCfg.GrpcClient)
	//if err != nil {
	//	return nil, err
	//}

	grokClient, err := grok_client.NewGrokClient(&appCfg.GrokCfg, proxyCfg)
	if err != nil {
		return nil, err
	}

	pgxPool, err := postgres.GetPool(context.Background(), appCfg.Database)
	if err != nil {
		return nil, err
	}

	postsRepo := postsRepoLib.NewPostsRepo(pgxPool)
	postsDelivery := posts_command.NewPostsHandler(postsRepo, nil, grokClient)

	return &PostsApp{
		appCfg: appCfg,
		grpc:   nil,
		http:   postsDelivery,
	}, nil
}
