package bots

import (
	"context"

	"github.com/google/uuid"
)

func (r *BotsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE bots SET is_deleted = true, updated_at = NOW() WHERE id=$1`, id)
	return err
}
