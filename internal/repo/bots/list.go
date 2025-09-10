package bots

import (
	"context"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *PGRepo) List(ctx context.Context) ([]model.Bot, error) {
	rows, err := r.pool.Query(ctx, `select id, bot_name, system_prompt from bots order by created_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Bot
	for rows.Next() {
		var b model.Bot
		if err = rows.Scan(&b.ID, &b.Name, &b.SystemPrompt); err != nil {
			return nil, err
		}
		res = append(res, b)
	}

	return res, rows.Err()
}
