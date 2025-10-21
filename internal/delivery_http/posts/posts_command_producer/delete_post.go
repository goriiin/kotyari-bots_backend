package posts_command_producer

import (
	"context"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
)

func (p *PostsCommandHandler) DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error) {
	// -> kafka

	//err := p.repo.DeletePost(ctx, uint64(params.PostId))
	//if err != nil {
	//	return &gen.DeletePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}

	return &gen.NoContent{}, nil
}
