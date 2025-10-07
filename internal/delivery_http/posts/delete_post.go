package posts

import (
	"context"
	"net/http"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts"
)

func (p *PostsHandler) DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error) {
	err := p.repo.DeletePost(ctx, uint64(params.PostId))
	if err != nil {
		return &gen.DeletePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return &gen.NoContent{}, nil
}
