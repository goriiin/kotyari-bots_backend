package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *Repository) Update(ctx context.Context, p model.Profile) error {
	_, err := r.db.Exec(ctx,
		`UPDATE profiles SET name=$2, email=$3, system_prompt=$4, updated_at=now() WHERE id=$1`,
		p.ID, p.Name, p.Email, p.SystemPromt)
	return err
}
