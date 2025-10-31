package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (p *PostsQueryRepo) ListPosts(ctx context.Context) ([]model.Post, error) {
	const query = `
		SELECT id, otveti_id, bot_id, profile_id, platform_type::text, post_type::text, post_title, post_text, created_at, updated_at
		FROM posts
	`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(constants.ErrInternal, "failed to query rows: %s", err.Error())
	}
	defer rows.Close()

	postsDTO, err := pgx.CollectRows(rows, pgx.RowToStructByName[posts.PostDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNotFound
		}

		return nil, errors.Wrapf(constants.ErrInternal, "failed to collect rows: %s", err)
	}

	return posts.PostsDTOToModel(postsDTO), nil
}
