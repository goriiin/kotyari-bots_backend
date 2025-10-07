package aggregator

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (a *AggregatorRepo) AddTopics(ctx context.Context, topics []model.Topic) error {
	if len(topics) == 0 {
		return errors.New("no topics to insert")
	}

	batch := &pgx.Batch{}
	for _, topic := range topics {
		batch.Queue(`
		INSERT INTO topics (source, text, hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (hash) DO NOTHING`,
			topic.Source, topic.Text, topic.Hash,
		)
	}

	b := a.db.SendBatch(ctx, batch)
	for range topics {
		_, err := b.Exec()
		if err != nil {
			return errors.Wrap(err, "failed to execute batch")
		}
	}

	return nil
}
