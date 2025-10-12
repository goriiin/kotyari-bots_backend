package bots

import (
	"context"

	"github.com/google/uuid"
)

func (r *BotsRepository) AddProfileID(ctx context.Context, botID, profileID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bots SET 
            profile_ids = array_append(COALESCE(profile_ids, '{}'), $2),
            profiles_count = COALESCE(array_length(array_append(profile_ids, $2), 1), 0),
            updated_at = now()
         WHERE id = $1 AND NOT ($2 = ANY(COALESCE(profile_ids, '{}')))`,
		botID, profileID)
	return err
}

func (r *BotsRepository) RemoveProfileID(ctx context.Context, botID, profileID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bots SET 
            profile_ids = array_remove(profile_ids, $2),
            profiles_count = COALESCE(array_length(array_remove(profile_ids, $2), 1), 0),
            updated_at = now()
         WHERE id = $1`,
		botID, profileID)
	return err
}
