package posts_query

import (
	"github.com/goriiin/kotyari-bots_backend/internal/repo/pool"
)

type PostsQueryRepo struct {
	db pool.DBPool
}

func NewPostsQueryRepo(dbPool pool.DBPool) *PostsQueryRepo {
	return &PostsQueryRepo{
		dbPool,
	}
}
