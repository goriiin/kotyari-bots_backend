package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) List(ctx context.Context) ([]model.Profile, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, name, email, system_prompt, created_at, updated_at 
			FROM profiles 
			WHERE user_id=$1 
			ORDER BY created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[profileDTO])
	if err != nil {
		return nil, err
	}

	return toModels(dtos), nil
}
