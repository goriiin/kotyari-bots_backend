package posts

import "github.com/goriiin/kotyari-bots_backend/internal/repo/pool"

type PostsRepo struct {
	// TODO: add logs
	db pool.DBPool
}

func NewPostsRepo(dbPool pool.DBPool) *PostsRepo {
	return &PostsRepo{
		dbPool,
	}
}
