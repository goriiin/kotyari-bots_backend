package profiles

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Profile, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return model.Profile{}, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, name, email, system_prompt, created_at, updated_at 
			FROM profiles 
			WHERE id=$1 AND user_id=$2`,
		id, userID)
	if err != nil {
		return model.Profile{}, err
	}

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[profileDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Profile{}, constants.ErrNotFound
		}
		return model.Profile{}, err
	}
	return dto.toModel(), nil
}

// GetByIDs is an internal service-to-service lookup (called over gRPC by the
// bots service to resolve a bot's profiles). The gRPC contract carries no
// user_id, so this is intentionally not tenant-scoped; callers must only pass
// IDs they are already authorized to read. User-facing reads go through
// GetByID, which is scoped by user_id.
func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, email, system_prompt, created_at, updated_at FROM profiles WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[profileDTO])
	if err != nil {
		return nil, err
	}

	return toModels(dtos), nil
}
