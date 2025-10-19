package posts_query

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/posts"
	"github.com/jackc/pgx/v5"
)

func (p *PostsQueryRepo) GetWithCategories(ctx context.Context, id uint64) (model.PostWithCategories, error) {
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
		return model.PostWithCategories{}, errors.Wrap(err, "failed to query row")
	}
	defer rows.Close()

	dtoCategories, err := pgx.CollectRows(rows, pgx.RowToStructByName[posts.CategoryDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// TODO: Сделать нормальную работу с ошибками
			return model.PostWithCategories{}, errors.New("failed to query categories")
		}

		return model.PostWithCategories{}, errors.Wrap(err, "failed to collect rows")
	}

	return model.PostWithCategories{
		Post:       post,
		Categories: posts.CategoriesDtoToModel(dtoCategories),
	}, nil
}
