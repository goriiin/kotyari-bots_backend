package posts_command

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (p *PostsCommandRepo) UpdatePostsBatch(ctx context.Context, posts []model.Post) (err error) {
	const query = `
        UPDATE posts
        SET post_title = $1,
            post_text = $2,
            updated_at = NOW()
        WHERE id = $3 AND group_id = $4
    `

	batch := &pgx.Batch{}

	for _, post := range posts {
		batch.Queue(query,
			post.Text,
			post.Title,
			post.ID,
			post.GroupID,
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
		if _, err = br.Exec(); err != nil {
			return errors.Wrapf(constants.ErrInternal, "error happened while inserting posts: %s", err.Error())
		}
	}

	return nil
}
