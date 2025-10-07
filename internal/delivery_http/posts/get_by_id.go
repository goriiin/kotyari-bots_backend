package posts

import (
	"context"
	"net/http"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts"
)

func (p *PostsHandler) GetPostById(ctx context.Context, params gen.GetPostByIdParams) (gen.GetPostByIdRes, error) {
	// TODO: По идее тут должна быть бизнес логика выбора какой пост фетчить (с категориями или без) в зависимости от платформы,
	// Пока возвращается пост без категорий
	post, err := p.repo.GetByID(ctx, uint64(params.PostId))
	if err != nil {
		return &gen.GetPostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return modelToHttp(post), nil
}
