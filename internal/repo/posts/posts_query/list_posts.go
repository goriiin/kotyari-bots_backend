package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	"github.com/jackc/pgx/v5"
)

func (p *PostsQueryRepo) ListPosts(ctx context.Context) ([]model.Post, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return nil, err
	}

	const query = `
		SELECT id, otveti_id, group_id, user_prompt, bot_id, bot_name, profile_id, profile_name, 
		       platform_type::text, post_type::text, post_title, post_text, created_at, updated_at
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrapf(constants.ErrInternal, "failed to query rows: %s", err.Error())
	}
	defer rows.Close()

	postsDTO, err := pgx.CollectRows(rows, pgx.RowToStructByName[posts.PostDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Empty list is not an error usually, but strict mapping might behave differently
			return []model.Post{}, nil
		}
		return nil, errors.Wrapf(constants.ErrInternal, "failed to collect rows: %s", err)
	}

	return posts.PostsDTOToModel(postsDTO), nil
}
