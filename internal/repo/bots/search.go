package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	"github.com/jackc/pgx/v5"
)

func (r BotsRepository) Search(ctx context.Context, query string) ([]model.Bot, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT 
		    id,
		    bot_name, 
		    system_prompt, 
		    moderation_required, 
		    profile_ids, 
		    profiles_count, 
		    created_at, 
		    updated_at
		FROM bots
		WHERE is_deleted = false
		  AND user_id = $2
		  AND (bot_name ILIKE $1 OR system_prompt ILIKE $1)
		ORDER BY created_at DESC
	`, query, userID)
	if err != nil {
		return nil, err
	}
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[botDTO])
	if err != nil {
		return nil, err
	}
	return toModels(dtos), rows.Err()
}
