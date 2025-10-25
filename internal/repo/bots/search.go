package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r BotsRepository) Search(ctx context.Context, query string) ([]model.Bot, error) {
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
		WHERE isdeleted = false
		  AND (botname ILIKE $1 OR systemprompt ILIKE $1)
		ORDER BY createdat DESC
	`, query)
	if err != nil {
		return nil, err
	}
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[botDTO])
	if err != nil {
		return nil, err
	}
	return toModels(dtos), rows.Err()
}
