package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *PGRepo) Update(ctx context.Context, b model.Bot) error {
	_, err := r.pool.Exec(ctx,
		`update bots set bot_name=$2, system_prompt=$3, moderation_required=$4,profile_ids=$5, updated_at=now() where id=$1`,
		b.ID, b.Name, b.SystemPrompt, b.ModerationRequired, b.ProfileIDs)
	if err != nil {
		return err
	}
	return nil
}
