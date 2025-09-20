package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *BotsRepository) Get(ctx context.Context, id uuid.UUID) (model.Bot, error) {
	var out model.Bot
	err := r.db.QueryRow(ctx, `select id, bot_name, system_prompt, moderation_required, profile_ids,auto_publish from bots where id=$1`, id).
		Scan(&out.ID, &out.Name, &out.SystemPrompt, &out.ModerationRequired, &out.ProfileIDs, &out.AutoPublish)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Bot{}, constants.ErrNotFound
		}
		return model.Bot{}, err
	}
	return out, nil
}
