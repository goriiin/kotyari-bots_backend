package bots

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	"github.com/jackc/pgx/v5"
)

func (r BotsRepository) Get(ctx context.Context, id uuid.UUID) (model.Bot, error) {
	// gRPC вызовы от других сервисов могут не иметь контекста пользователя.
	// Если user_id нет, делаем запрос только по ID (внутреннее доверие).
	// Если user_id есть (HTTP API), фильтруем по нему.

	if userID, err := user.GetID(ctx); err == nil {
		rows, err := r.db.Query(ctx, `
			SELECT 
				id, bot_name, system_prompt, moderation_required, 
				profile_ids, profiles_count, created_at, updated_at
			FROM bots
			WHERE id = $1 AND is_deleted = false AND user_id = $2
		`, id, userID)
		if err != nil {
			return model.Bot{}, err
		}
		dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[botDTO])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.Bot{}, constants.ErrNotFound
			}
			return model.Bot{}, err
		}
		return dto.toModel(), nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT 
		    id, bot_name, system_prompt, moderation_required, 
		    profile_ids, profiles_count, created_at, updated_at
		FROM bots
		WHERE id = $1 AND is_deleted = false
	`, id)
	if err != nil {
		return model.Bot{}, err
	}
	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[botDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Bot{}, constants.ErrNotFound
		}
		return model.Bot{}, err
	}
	return dto.toModel(), nil
}
