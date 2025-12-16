package posts_query

import (
	"context"
	"net/http"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
)

func (p *PostsQueryHandler) GetPostById(ctx context.Context, params gen.GetPostByIdParams) (gen.GetPostByIdRes, error) {
	post, err := p.repo.GetByID(ctx, params.PostId)
	if err != nil {
		p.log.Error(err, true, "failed to get post by id")
		return &gen.GetPostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return posts.QueryModelToHttp(post), nil
}
