package bots

import (
	"context"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *BotsRepository) GetSummary(ctx context.Context) (model.BotsSummary, error) {
	var summary model.BotsSummary
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(profiles_count), 0) 
         FROM bots 
         WHERE is_deleted = false`).
		Scan(&summary.TotalBots, &summary.TotalProfilesAttached)
	if err != nil {
		return model.BotsSummary{}, err
	}
	return summary, nil
}
