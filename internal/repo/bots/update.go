package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
)

func (r *BotsRepository) Update(ctx context.Context, b model.Bot) error {
	userID, err := user.GetID(ctx)
	if err != nil {
		return err
	}

	tag, err := r.db.Exec(ctx,
		`UPDATE bots SET 
                bot_name = $2, 
                system_prompt = $3,
                moderation_required = $4, 
                profile_ids = $5,
                profiles_count = $6, 
                updated_at = now()
            WHERE id = $1 AND user_id = $7`,
		b.ID, b.Name, b.SystemPrompt, b.ModerationRequired, b.ProfileIDs, len(b.ProfileIDs), userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return constants.ErrNotFound
	}
	return nil
}
