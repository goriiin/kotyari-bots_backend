package posts_command

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *PostsCommandRepo) CreatePost(ctx context.Context, post model.Post, categoryIDs []uuid.UUID) (model.Post, error) {
	txOpts := pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	}

	tx, err := p.db.BeginTx(ctx, txOpts)
	if err != nil {
		return model.Post{}, errors.Wrapf(constants.ErrInternal, "failed to begin transaction: %s", err.Error())
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			// TODO: log err
		}
	}()

	var postType pgtype.Text
	if post.Type != "" {
		postType = pgtype.Text{String: string(post.Type), Valid: true}
	}

	const query = `
		INSERT INTO posts (id, otveti_id, bot_id, profile_id, group_id, user_prompt, platform_type, post_type, post_title, post_text)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		post.ID,
		post.OtvetiID,
		post.BotID,
		post.ProfileID,
		post.GroupID,
		post.UserPrompt,
		post.Platform,
		postType,
		post.Title,
		post.Text)
	if err = row.Scan(&post.CreatedAt, &post.UpdatedAt); err != nil {
		return model.Post{}, errors.Wrapf(constants.ErrInternal, "failed to scan row: %s", err.Error())
	}

	if len(categoryIDs) > 0 {
		err = p.insertPostCategoriesBatch(ctx, tx, post.ID, categoryIDs)
		return model.Post{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Post{}, errors.Wrap(constants.ErrInternal, "failed to commit transaction")
	}

	return post, nil
}

func (p *PostsCommandRepo) insertPostCategoriesBatch(ctx context.Context, tx pgx.Tx, postID uuid.UUID, categoryIDs []uuid.UUID) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	const query = `
		INSERT INTO post_categories (post_id, category_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	b := &pgx.Batch{}
	for _, cid := range categoryIDs {
		b.Queue(query, postID, cid)
	}

	br := tx.SendBatch(ctx, b)
	for i := 0; i < b.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return errors.Wrapf(constants.ErrInternal, "error happened while inserting categories: %s", err.Error())
		}
	}
	return br.Close()
}
