package bots

import (
	"context"

	"github.com/google/uuid"
)

func (r *PGRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `delete from bots where id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}
