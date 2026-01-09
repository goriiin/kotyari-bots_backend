package bots

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *BotsRepository) GetSummary(ctx context.Context) (model.BotsSummary, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return model.BotsSummary{}, err
	}

	var summary model.BotsSummary
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(profiles_count), 0) 
         FROM bots 
         WHERE is_deleted = false AND user_id = $1`, userID).
		Scan(&summary.TotalBots, &summary.TotalProfilesAttached)
	if err != nil {
		return model.BotsSummary{}, err
	}
	return summary, nil
}
