package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *Repository) List(ctx context.Context) ([]model.Profile, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, email, system_prompt, created_at, updated_at FROM profiles ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Profile
	for rows.Next() {
		var p model.Profile
		if err = rows.Scan(&p.ID, &p.Name, &p.Email, &p.SystemPromt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, rows.Err()
}
