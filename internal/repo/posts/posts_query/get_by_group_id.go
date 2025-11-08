package posts_query

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (p *PostsQueryRepo) GetByGroupId(ctx context.Context, groupID uuid.UUID) ([]model.Post, error) {
	fmt.Println(groupID)

	const query = `
		SELECT id, otveti_id, group_id, user_prompt, bot_id, bot_name, profile_id, profile_name, 
		       platform_type::text, post_type::text, post_title, post_text, created_at, updated_at
		FROM posts
		WHERE group_id = $1
	`

	rows, err := p.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, errors.Wrapf(constants.ErrInternal, "failed to query rows: %s", err.Error())
	}
	defer rows.Close()

	postsDTO, err := pgx.CollectRows(rows, pgx.RowToStructByName[posts.PostDTO])
	if err != nil {
		fmt.Println(postsDTO)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNotFound
		}

		return nil, errors.Wrapf(constants.ErrInternal, "failed to collect rows: %s", err)
	}

	return posts.PostsDTOToModel(postsDTO), nil
}
