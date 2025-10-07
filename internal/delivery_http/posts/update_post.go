package posts

import (
	"context"
	"net/http"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (p *PostsHandler) UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error) {
	oldPost := model.Post{
		ID:    uint64(params.PostId),
		Title: req.Title,
		Text:  req.Text,
	}

	modifiedPost, err := p.repo.UpdatePost(ctx, oldPost)
	if err != nil {
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return modelToHttp(modifiedPost), nil
}
