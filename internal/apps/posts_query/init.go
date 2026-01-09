package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	postsQueryHandler "github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_query"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	postsQueryRepo "github.com/goriiin/kotyari-bots_backend/internal/repo/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type postsGetter interface {
	GetPostById(ctx context.Context, params gen.GetPostByIdParams) (gen.GetPostByIdRes, error)
	ListPosts(ctx context.Context) (gen.ListPostsRes, error)
	CheckGroupId(ctx context.Context, params gen.CheckGroupIdParams) (gen.CheckGroupIdRes, error)
	CheckGroupIds(ctx context.Context) (gen.CheckGroupIdsRes, error)
}

type PostsQueryApp struct {
	handler postsGetter
	config  *PostsQueryConfig
}

func NewPostsQueryApp(config *PostsQueryConfig) (*PostsQueryApp, error) {
	log := logger.NewLogger("posts-query", &config.ConfigBase)

	pool, err := postgres.GetPool(context.Background(), config.Database)

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	repo := postsQueryRepo.NewPostsQueryRepo(pool)

	handler := postsQueryHandler.NewPostsQueryHandler(repo, log)

	return &PostsQueryApp{
		handler: handler,
		config:  config,
	}, nil
}
