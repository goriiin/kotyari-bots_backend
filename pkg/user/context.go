package user

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

type ctxKey struct{}

var userCtxKey = ctxKey{}

var ErrNoUser = errors.New("user not found in context")

func WithID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userCtxKey, id)
}

func GetID(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(userCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrNoUser
	}
	return id, nil
}
