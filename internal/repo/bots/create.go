package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/internal/pkg/user"
)

func (r BotsRepository) Create(ctx context.Context, b model.Bot) error {
	userID, err := user.GetID(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO bots (
		                  id, 
		                  bot_name,
		                  system_prompt,
		                  moderation_required,
		                  profile_ids,
		                  profiles_count,
		                  user_id)
		VALUES (
			$1, $2, $3, $4, $5::uuid[],
			COALESCE(array_length($5::uuid[], 1), 0),
			$6
		)
	`, b.ID, b.Name, b.SystemPrompt, b.ModerationRequired, b.ProfileIDs, userID)
	return err
}
