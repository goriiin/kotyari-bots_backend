package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	postsQueryHandler "github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_query"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	postsQueryRepo "github.com/goriiin/kotyari-bots_backend/internal/repo/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type postsGetter interface {
	GetPostById(ctx context.Context, params gen.GetPostByIdParams) (gen.GetPostByIdRes, error)
	ListPosts(ctx context.Context) (gen.ListPostsRes, error)
}

type PostsQueryApp struct {
	handler postsGetter
}

func NewPostsQueryApp() (*PostsQueryApp, error) {
	// TODO: вынести в конфиг
	pool, err := postgres.GetPool(context.Background(), postgres.Config{
		Host:     "posts_db",
		Port:     5432,
		Name:     "posts",
		User:     "postgres",
		Password: "123",
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	repo := postsQueryRepo.NewPostsQueryRepo(pool)

	handler := postsQueryHandler.NewPostsQueryHandler(repo)

	return &PostsQueryApp{handler: handler}, nil
}
