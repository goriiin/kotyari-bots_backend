package posts_command

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
)

func (p *PostsCommandRepo) CreatePostsBatch(ctx context.Context, posts []model.Post) (err error) {
	const query = `
		INSERT INTO posts (id, otveti_id, bot_id, bot_name, profile_id, profile_name, group_id, user_prompt, platform_type, post_type, post_title, post_text)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING created_at, updated_at
	`

	batch := &pgx.Batch{}

	for _, post := range posts {
		batch.Queue(query,
			post.ID,
			post.OtvetiID,
			post.BotID,
			post.BotName,
			post.ProfileID,
			post.ProfileName,
			post.GroupID,
			post.UserPrompt,
			post.Platform,
			post.Type,
			post.Title,
			post.Text,
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
