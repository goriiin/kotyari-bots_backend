package posts_command

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/jackc/pgx/v5"
)

func (p *PostsCommandRepo) UpdatePost(ctx context.Context, post model.Post) (model.Post, error) {
	const query = `
		UPDATE posts
		SET post_title=$1, post_text=$2, updated_at=NOW()
		WHERE id=$3
		RETURNING id, bot_id, profile_id, platform_type, post_type, post_title, post_text, created_at, updated_at
	`

	rows, err := p.db.Query(ctx, query, post.Title, post.Text, post.ID)
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to update post")
	}

	modifiedPost, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[posts.PostDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Post{}, errors.Wrap(err, "post not found")
		}

		return model.Post{}, errors.Wrap(err, "unexpected error happened")
	}
	return modifiedPost.ToModel(), nil
}
