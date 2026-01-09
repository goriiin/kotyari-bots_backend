package profiles

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
)

func (r *Repository) Create(ctx context.Context, p model.Profile) error {
	userID, err := user.GetID(ctx)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx,
		`INSERT INTO profiles (id, name, email, system_prompt, user_id) VALUES ($1, $2, $3, $4, $5)`,
		p.ID, p.Name, p.Email, p.SystemPromt, userID)
	return err
}
