package posts_query

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (p *PostsQueryHandler) CheckGroupId(ctx context.Context, params gen.CheckGroupIdParams) (gen.CheckGroupIdRes, error) {
	fmt.Println(params)

	groupPosts, err := p.repo.GetByGroupId(ctx, params.GroupId)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), constants.NotFoundMsg):
			return &gen.CheckGroupIdNotFound{
				ErrorCode: http.StatusNotFound,
				Message:   "Посты с этим groupID еще не готовы",
			}, nil

		case strings.Contains(err.Error(), constants.InternalMsg):
			return &gen.CheckGroupIdInternalServerError{
				ErrorCode: http.StatusInternalServerError,
				Message:   err.Error(),
			}, nil
		}
	}

	return posts.QueryPostsToHttp(groupPosts), nil
}
