package profiles

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM profiles WHERE id=$1`, id)
	return err
}
