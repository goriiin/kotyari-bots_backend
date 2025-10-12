package profiles

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Profile, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, email, system_prompt, created_at, updated_at FROM profiles WHERE id=$1`, id)
	if err != nil {
		return model.Profile{}, err
	}

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[profileDTO])
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Profile{}, constants.ErrNotFound
		}
		return model.Profile{}, err
	}
	return dto.toModel(), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, email, system_prompt, created_at, updated_at FROM profiles WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[profileDTO])
	if err != nil {
		return nil, err
	}

	return toModels(dtos), nil
}
