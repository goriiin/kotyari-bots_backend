package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
)

func (r *Repository) Update(ctx context.Context, p model.Profile) error {
	userID, err := user.GetID(ctx)
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx,
		`UPDATE profiles SET name=$2, email=$3, system_prompt=$4, updated_at=now()
                WHERE id=$1 AND user_id=$5`,
		p.ID, p.Name, p.Email, p.SystemPromt, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return constants.ErrNotFound
	}
	return nil
}
