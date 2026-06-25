package posts_query

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (p *PostsQueryHandler) CheckGroupId(ctx context.Context, params gen.CheckGroupIdParams) (gen.CheckGroupIdRes, error) {
	groupPosts, err := p.repo.GetByGroupId(ctx, params.GroupId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.CheckGroupIdNotFound{
				ErrorCode: http.StatusNotFound,
				Message:   "Посты с этим groupID еще не готовы",
			}, nil
		}

		p.log.Error(err, true, "CheckGroupId: get by group id")
		return &gen.CheckGroupIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	return posts.QueryPostsToHttp(groupPosts), nil
}
