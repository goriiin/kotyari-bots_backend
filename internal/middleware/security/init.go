package security

import (
	"context"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
)

type NotImplemented struct{}

func (n *NotImplemented) HandleCookieAuth(ctx context.Context, _ gen.OperationName, _ gen.CookieAuth) (context.Context, error) {
	return ctx, nil
}

func (n *NotImplemented) HandleCsrfAuth(ctx context.Context, _ gen.OperationName, _ gen.CsrfAuth) (context.Context, error) {
	return ctx, nil
}
