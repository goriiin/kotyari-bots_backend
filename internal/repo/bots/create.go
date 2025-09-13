package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *BotsRepository) Create(ctx context.Context, b model.Bot) error {
	_, err := r.db.Exec(ctx,
		`insert into bots (id, bot_name, system_prompt, moderation_required, profile_ids) values ($1, $2, $3, $4)`,
		b.ID, b.Name, b.SystemPrompt, b.ModerationRequired, b.ProfileIDs)
	return err
}
