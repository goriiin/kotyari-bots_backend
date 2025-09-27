package profiles

import (
	"context"
	"github.com/google/uuid"
)

// Реализация запроса к БД для проверки существования профилей
func (r *Repository) Exist(ctx context.Context, ids []uuid.UUID) (map[string]bool, error) {
	rows, err := r.pool.Query(ctx, `SELECT id FROM profiles WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existenceMap := make(map[string]bool, len(ids))
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existenceMap[id.String()] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return existenceMap, nil
}
