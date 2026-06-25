package posts_command

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
)

func (p *PostsCommandRepo) DeletePost(ctx context.Context, id, userID uuid.UUID) error {
	const query = `
		DELETE FROM posts WHERE id=$1 AND user_id=$2
	`

	ct, err := p.db.Exec(ctx, query, id, userID)
	if err != nil {
		return errors.Wrapf(constants.ErrInternal, "failed to delete post: %s", err.Error())
	}

	if ct.RowsAffected() == 0 {
		return errors.Wrap(constants.ErrNotFound, "post not found")
	}

	return nil
}
