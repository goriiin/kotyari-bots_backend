package posts_query

import (
	"context"
	"net/http"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (p *PostsQueryHandler) CheckGroupIds(ctx context.Context) (gen.CheckGroupIdsRes, error) {
	postsStatuses, err := p.repo.CheckGroupIds(ctx)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), constants.NotFoundMsg):
			return &gen.CheckGroupIdsNotFound{
				ErrorCode: http.StatusNotFound,
				Message:   "Постов нет",
			}, nil

		case strings.Contains(err.Error(), constants.InternalMsg):
			return &gen.CheckGroupIdsNotFound{
				ErrorCode: http.StatusInternalServerError,
				Message:   err.Error(),
			}, nil
		}
	}

	return posts.PostsCheckModelsToHttpSlice(postsStatuses), nil
}
