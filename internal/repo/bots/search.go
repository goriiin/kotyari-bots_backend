package bots

import (
	"context"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *BotsRepository) Search(ctx context.Context, query string) ([]model.Bot, error) {
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		return []model.Bot{}, nil
	}
	rows, err := r.db.Query(ctx,
		`SELECT id, bot_name, system_prompt, moderation_required, profile_ids, profiles_count, created_at, updated_at 
         FROM bots 
         WHERE (bot_name ILIKE $1 OR system_prompt ILIKE $1) AND is_deleted = false
         ORDER BY created_at DESC`,
		"%"+trimmedQuery+"%")
	if err != nil {
		return nil, err
	}

	bots, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Bot])
	if err != nil {
		return nil, err
	}
	return bots, nil
}
