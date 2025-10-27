package posts_command

import "github.com/goriiin/kotyari-bots_backend/internal/repo/pool"

// PostsCommandRepo TODO: Добавить функцию добавления постов батчами
type PostsCommandRepo struct {
	db pool.DBPool
}

func NewPostsCommandRepo(dbPool pool.DBPool) *PostsCommandRepo {
	return &PostsCommandRepo{
		dbPool,
	}
}
