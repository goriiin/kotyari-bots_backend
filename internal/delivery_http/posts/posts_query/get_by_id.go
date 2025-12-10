package posts_query

import (
	"context"
	"net/http"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (p *PostsQueryHandler) GetPostById(ctx context.Context, params gen.GetPostByIdParams) (gen.GetPostByIdRes, error) {
	// TODO: По идее тут должна быть бизнес логика выбора какой пост фетчить (с категориями или без) в зависимости от платформы,
	// Пока возвращается пост без категорий
	post, err := p.repo.GetByID(ctx, params.PostId)
	if err != nil {

		if strings.Contains(err.Error(), constants.NotFoundMsg) {
			return &gen.GetPostByIdNotFound{ErrorCode: http.StatusNotFound, Message: "post not found"}, nil
		}

		return &gen.GetPostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return posts.QueryModelToHttp(post), nil
}
