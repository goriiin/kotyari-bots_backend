package bots

import (
	"context"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (r BotsRepository) AddProfileID(ctx context.Context, botID, profileID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE bots
		SET
			profile_ids    = array_append(COALESCE(profile_ids, '{}'::uuid[]), $2),
			profiles_count = COALESCE(array_length(array_append(COALESCE(profile_ids, '{}'::uuid[]), $2), 1), 0),
			updated_at     = now()
		WHERE id = $1
		  AND NOT $2 = ANY(COALESCE(profile_ids, '{}'::uuid[]))
	`, botID, profileID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return constants.ErrNotFound
	}
	return nil
}

func (r BotsRepository) RemoveProfileID(ctx context.Context, botID, profileID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE bots
		SET
			profile_ids    = array_remove(COALESCE(profile_ids, '{}'::uuid[]), $2),
			profiles_count = COALESCE(array_length(array_remove(COALESCE(profile_ids, '{}'::uuid[]), $2), 1), 0),
			updated_at     = now()
		WHERE id = $1
	`, botID, profileID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return constants.ErrNotFound
	}
	return nil
}
