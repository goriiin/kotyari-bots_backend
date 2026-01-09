package posts_query

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type postsQueryRepository interface {
	GetWithCategories(ctx context.Context, id uuid.UUID) (model.PostWithCategories, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Post, error)
	ListPosts(ctx context.Context) ([]model.Post, error)
	GetByGroupId(ctx context.Context, groupID uuid.UUID) ([]model.Post, error)
	CheckGroupIds(ctx context.Context) ([]model.Post, error)
}

type PostsQueryHandler struct {
	repo postsQueryRepository
	log  *logger.Logger
}

func NewPostsQueryHandler(repo postsQueryRepository, log *logger.Logger) *PostsQueryHandler {
	return &PostsQueryHandler{
		repo: repo,
		log:  log,
	}
}
