package bots

import (
	"context"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	"github.com/jackc/pgx/v5"
)

// likeEscaper escapes the LIKE/ILIKE wildcard characters so that a user-supplied
// search term is matched literally as a substring instead of being interpreted
// as a pattern. The backslash is used as the ESCAPE character in the query.
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func (r BotsRepository) Search(ctx context.Context, query string) ([]model.Bot, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return nil, err
	}

	// Wrap the (escaped) term in wildcards so ILIKE performs a substring match
	// rather than an exact, case-insensitive comparison.
	pattern := "%" + likeEscaper.Replace(query) + "%"

	rows, err := r.db.Query(ctx, `
		SELECT 
		    id,
		    bot_name, 
		    system_prompt, 
		    moderation_required, 
		    profile_ids, 
		    profiles_count, 
		    created_at, 
		    updated_at
		FROM bots
		WHERE is_deleted = false
		  AND user_id = $2
		  AND (bot_name ILIKE $1 ESCAPE '\' OR system_prompt ILIKE $1 ESCAPE '\')
		ORDER BY created_at DESC
	`, pattern, userID)
	if err != nil {
		return nil, err
	}
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[botDTO])
	if err != nil {
		return nil, err
	}
	return toModels(dtos), rows.Err()
}
