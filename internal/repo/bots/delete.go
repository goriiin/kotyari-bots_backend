package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
)

func (r *BotsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	userID, err := user.GetID(ctx)
	if err != nil {
		return err
	}

	tag, err := r.db.Exec(ctx, `UPDATE bots SET is_deleted = true, updated_at = NOW() WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return constants.ErrNotFound
	}
	return nil
}
