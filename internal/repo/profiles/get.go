package profiles

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Profile, error) {
	var out model.Profile
	err := r.pool.QueryRow(ctx, `SELECT id, name, email, system_prompt, created_at, updated_at FROM profiles WHERE id=$1`, id).
		Scan(&out.ID, &out.Name, &out.Email, &out.SystemPromt, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Profile{}, constants.ErrNotFound
		}
		return model.Profile{}, err
	}
	return out, nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, email, system_prompt, created_at, updated_at FROM profiles WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []model.Profile
	for rows.Next() {
		var p model.Profile
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.SystemPromt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	return profiles, rows.Err()
}
