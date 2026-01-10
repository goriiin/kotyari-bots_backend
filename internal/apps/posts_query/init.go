package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/adapters/auth"
	postsQueryHandler "github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_query"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	postsQueryRepo "github.com/goriiin/kotyari-bots_backend/internal/repo/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
)

type postsGetter interface {
	GetPostById(ctx context.Context, params gen.GetPostByIdParams) (gen.GetPostByIdRes, error)
	ListPosts(ctx context.Context) (gen.ListPostsRes, error)
	CheckGroupId(ctx context.Context, params gen.CheckGroupIdParams) (gen.CheckGroupIdRes, error)
	CheckGroupIds(ctx context.Context) (gen.CheckGroupIdsRes, error)
}

type PostsQueryApp struct {
	handler    postsGetter
	config     *PostsQueryConfig
	authClient *auth.Client
}

type securityHandler struct {
	authClient *auth.Client
}

func (s *securityHandler) HandleSessionAuth(ctx context.Context, operationName string, t gen.SessionAuth) (context.Context, error) {
	userID, err := s.authClient.VerifySession(ctx, t.APIKey)
	if err != nil {
		return nil, err
	}
	return user.WithID(ctx, userID), nil
}

func NewPostsQueryApp(config *PostsQueryConfig) (*PostsQueryApp, error) {
	log := logger.NewLogger("posts-query", &config.ConfigBase)

	pool, err := postgres.GetPool(context.Background(), config.Database)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	authClient, err := auth.NewClient(config.Auth, log)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth client")
	}

	repo := postsQueryRepo.NewPostsQueryRepo(pool)
	handler := postsQueryHandler.NewPostsQueryHandler(repo, log)

	return &PostsQueryApp{
		handler:    handler,
		config:     config,
		authClient: authClient,
	}, nil
}
