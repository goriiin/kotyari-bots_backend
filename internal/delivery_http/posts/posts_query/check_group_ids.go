package posts_query

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (p *PostsQueryHandler) CheckGroupIds(ctx context.Context) (gen.CheckGroupIdsRes, error) {
	postsStatuses, err := p.repo.CheckGroupIds(ctx)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.CheckGroupIdsNotFound{
				ErrorCode: http.StatusNotFound,
				Message:   "Постов нет",
			}, nil
		}

		p.log.Error(err, true, "CheckGroupIds: check group ids")
		return &gen.CheckGroupIdsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	return posts.PostsCheckModelsToHttpSlice(postsStatuses), nil
}
