package posts_query

import (
	"context"
	"net/http"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
)

func (p *PostsQueryHandler) ListPosts(ctx context.Context) (gen.ListPostsRes, error) {
	postsModels, err := p.repo.ListPosts(ctx)
	if err != nil {
		return &gen.ListPostsInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return posts.QueryPostsToHttp(postsModels), nil
}
