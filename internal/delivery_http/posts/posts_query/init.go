package posts_query

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type postsQueryRepository interface {
	GetWithCategories(ctx context.Context, id uint64) (model.PostWithCategories, error)
	GetByID(ctx context.Context, id uint64) (model.Post, error)
}

type PostsQueryHandler struct {
	repo postsQueryRepository
}

func NewPostsQueryHandler(repo postsQueryRepository) *PostsQueryHandler {
	return &PostsQueryHandler{
		repo: repo,
	}
}
