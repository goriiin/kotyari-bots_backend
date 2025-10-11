package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *BotsRepository) Update(ctx context.Context, b model.Bot) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bots SET 
            bot_name=$2, 
            system_prompt=$3, 
            moderation_required=$4, 
            updated_at=now() 
         WHERE id=$1`,
		b.ID, b.Name, b.SystemPrompt, b.ModerationRequired)
	return err
}
