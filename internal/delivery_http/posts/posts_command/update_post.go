package posts_command

import (
	"context"
	"net/http"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (p *PostsCommandHandler) UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error) {
	oldPost := model.Post{
		ID:    uint64(params.PostId),
		Title: req.Title,
		Text:  req.Text,
	}

	// -> kafka

	//modifiedPost, err := p.repo.UpdatePost(ctx, oldPost)
	//if err != nil {
	//	return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}

	return posts.ModelToHttp(modifiedPost), nil
}
