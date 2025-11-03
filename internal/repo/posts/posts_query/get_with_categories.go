package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (p *PostsQueryRepo) GetWithCategories(ctx context.Context, id uuid.UUID) (model.PostWithCategories, error) {
	post, err := p.GetByID(ctx, id)
	if err != nil {
		return model.PostWithCategories{}, err
	}

	const query = `
		SELECT c.id, c.category_name, c.created_at, c.updated_at
		FROM categories c
		JOIN post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = $1
	`

	rows, err := p.db.Query(ctx, query, id)
	if err != nil {
		return model.PostWithCategories{}, errors.Wrapf(constants.ErrInternal, "failed to query row: %s", err.Error())
	}
	defer rows.Close()

	dtoCategories, err := pgx.CollectRows(rows, pgx.RowToStructByName[posts.CategoryDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PostWithCategories{}, constants.ErrNotFound
		}

		return model.PostWithCategories{}, errors.Wrapf(constants.ErrInternal, "failed to collect rows: %s", err.Error())
	}

	return model.PostWithCategories{
		Post:       post,
		Categories: posts.CategoriesDtoToModel(dtoCategories),
	}, nil
}
