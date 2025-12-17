package posts_command

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (p *PostsCommandRepo) SeenPostsBatch(ctx context.Context, postsIds []uuid.UUID) (err error) {
	const query = `
        UPDATE posts
        SET is_seen = $1
        WHERE id = $2
    `

	batch := &pgx.Batch{}

	for _, id := range postsIds {
		batch.Queue(query,
			true,
			id,
		)
	}

	br := p.db.SendBatch(ctx, batch)
	defer func() {
		err = br.Close()
		if err != nil {
			fmt.Printf("Error happened closing batch: %s\n", err.Error())
		}
	}()

	for i := 0; i < batch.Len(); i++ {
		ct, err := br.Exec()
		if err != nil {
			return errors.Wrapf(constants.ErrInternal, "error happened while inserting posts: %s", err.Error())
		}

		// TODO: NOT WORKING
		if ct.RowsAffected() == 0 {
			return errors.Wrapf(constants.ErrNotFound, "posts not found")
		}
	}

	return nil
}
