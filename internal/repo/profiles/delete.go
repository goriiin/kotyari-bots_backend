package profiles

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	userID, err := user.GetID(ctx)
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, `DELETE FROM profiles WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return constants.ErrNotFound
	}
	return nil
}
