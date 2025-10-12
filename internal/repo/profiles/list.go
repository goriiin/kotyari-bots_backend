package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) List(ctx context.Context) ([]model.Profile, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, email, system_prompt, created_at, updated_at FROM profiles ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[profileDTO])
	if err != nil {
		return nil, err
	}

	return toModels(dtos), nil
}
