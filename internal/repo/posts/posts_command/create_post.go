package posts_command

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *PostsCommandRepo) CreatePost(ctx context.Context, post model.Post, categoryIDs []uuid.UUID) (model.Post, error) {
	txOpts := pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	}

	tx, err := p.db.BeginTx(ctx, txOpts)
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to begin transaction")
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
		INSERT INTO posts (bot_id, profile_id, platform_type, post_type, post_title, post_text)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query, post.BotID, post.ProfileID, postType, post.Title, post.Text)

	var postID uint64
	if err = row.Scan(&postID, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return model.Post{}, errors.Wrap(err, "failed to scan row")
	}
	post.ID = postID

	if len(categoryIDs) > 0 {
		err = p.insertPostCategoriesBatch(ctx, tx, postID, categoryIDs)
		return model.Post{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to commit transaction")
	}

	return post, nil
}

func (p *PostsCommandRepo) insertPostCategoriesBatch(ctx context.Context, tx pgx.Tx, postID uint64, categoryIDs []uuid.UUID) error {
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
			return errors.Wrap(err, "error happened while inserting categories")
		}
	}
	return br.Close()
}
