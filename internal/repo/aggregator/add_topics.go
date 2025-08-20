package aggregator

import (
	"context"
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
)

func (a *AggregatorRepo) AddTopics(ctx context.Context, topics []model.Topic) error {
	if len(topics) == 0 {
		return fmt.Errorf("no topics to insert")
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
			return fmt.Errorf("failed to execute batch, %w", err)
		}
	}

	return nil
}
