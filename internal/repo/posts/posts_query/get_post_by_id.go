package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/jackc/pgx/v5"
)

func (p *PostsQueryRepo) GetByID(ctx context.Context, id uuid.UUID) (model.Post, error) {
	const query = `
		SELECT id, otveti_id, bot_id, profile_id, platform_type::text, post_type::text, post_title, post_text, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	rows, err := p.db.Query(ctx, query, id)
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to query rows")
	}
	defer rows.Close()

	post, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[posts.PostDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// TODO: Сделать нормальную работу с ошибками
			return model.Post{}, errors.New("no such post")
		}

		return model.Post{}, errors.Wrap(err, "failed to collect rows")
	}

	return post.ToModel(), nil
}
