package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *BotsRepository) List(ctx context.Context) ([]model.Bot, error) {
	rows, err := r.db.Query(
		ctx,
		`
			select id, bot_name, system_prompt, profiles_count, created_at, updated_at
			from bots 
			order by created_at desc
			`,
	)
	if err != nil {
		return nil, err
	}

	bots, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Bot])
	if err != nil {
		return nil, err
	}

	return bots, nil
}
