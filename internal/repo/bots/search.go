package bots

import (
	"context"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *BotsRepository) Search(ctx context.Context, query string) ([]model.Bot, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, bot_name, system_prompt, moderation_required, profile_ids, profiles_count, created_at, updated_at 
         FROM bots 
         WHERE (bot_name ILIKE $1 OR system_prompt ILIKE $1) AND is_deleted = false
         ORDER BY created_at DESC`,
		"%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []model.Bot
	for rows.Next() {
		var b model.Bot
		if err := rows.Scan(&b.ID, &b.Name, &b.SystemPrompt, &b.ModerationRequired, &b.ProfileIDs, &b.ProfilesCount, &b.CreatedAt, &b.UpdateAt); err != nil {
			return nil, err
		}
		bots = append(bots, b)
	}
	return bots, rows.Err()
}
