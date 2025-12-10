package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (p *PostsQueryRepo) CheckGroupIds(ctx context.Context) ([]model.Post, error) {
	const query = `
		SELECT id, group_id, post_title, post_text, is_seen
		FROM posts
	`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(constants.ErrInternal, "failed to query rows: %s", err.Error())
	}
	defer rows.Close()

	postCheckDTO, err := pgx.CollectRows(rows, pgx.RowToStructByName[posts.PostCheckDTO])
	if err != nil {
		return nil, errors.Wrapf(constants.ErrInternal, "failed to collect rows: %s", err)
	}

	return posts.PostCheckDTOToModelSlice(postCheckDTO), nil
}
