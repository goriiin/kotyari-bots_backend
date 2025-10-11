package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (r *BotsRepository) Get(ctx context.Context, id uuid.UUID) (model.Bot, error) {
	rows, err := r.db.Query(ctx,
		`
			select id, bot_name, system_prompt, moderation_required, profile_ids, profiles_count, created_at, updated_at 
			from bots where id=$1 and is_deleted = false
			`,
		id)
	if err != nil {
		return model.Bot{}, err
	}

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[botDTO])
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Bot{}, constants.ErrNotFound
		}
		return model.Bot{}, err
	}

	return dto.toModel(), nil
}
