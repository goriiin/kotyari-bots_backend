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
			SELECT id, bot_name, system_prompt, moderation_required, profile_ids, profiles_count, created_at, updated_at
			FROM bots 
			WHERE is_deleted = false
			ORDER BY created_at DESC
			`,
	)
	if err != nil {
		return nil, err
	}

	botDTOs, err := pgx.CollectRows(rows, pgx.RowToStructByName[botDTO])
	if err != nil {
		return nil, err
	}

	return toModels(botDTOs), nil
}
