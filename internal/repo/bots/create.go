package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r BotsRepository) Create(ctx context.Context, b model.Bot) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO bots (
		                  id, 
		                  bot_name,
		                  system_prompt,
		                  moderation_required,
		                  profile_ids,
		                  profiles_count)
		VALUES (
			$1, $2, $3, $4, $5::uuid[],
			COALESCE(array_length($5::uuid[], 1), 0)
		)
	`, b.ID, b.Name, b.SystemPrompt, b.ModerationRequired, b.ProfileIDs)
	return err
}
