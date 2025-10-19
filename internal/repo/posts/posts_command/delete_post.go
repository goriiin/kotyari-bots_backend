package posts_command

import (
	"context"

	"github.com/go-faster/errors"
)

func (p *PostsCommandRepo) DeletePost(ctx context.Context, id uint64) error {
	const query = `
		DELETE FROM posts WHERE id=$1
	`

	ct, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "unexpected error happened")
	}

	if ct.RowsAffected() == 0 {
		// TODO: errors
		return errors.New("no rows affected")
	}

	return nil
}
