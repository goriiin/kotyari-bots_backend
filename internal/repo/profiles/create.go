package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (r *Repository) Create(ctx context.Context, p model.Profile) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO profiles (id, name, email, system_prompt) VALUES ($1, $2, $3, $4)`,
		p.ID, p.Name, p.Email, p.SystemPromt)
	return err
}
