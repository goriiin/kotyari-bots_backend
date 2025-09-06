package aggregator

import (
	"context"
)

func (a *AggregatorDelivery) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messageBatches := a.consumer.ReadBatches(ctx)

	for batch := range messageBatches {
		err := a.manager.AddTopics(ctx, batch)
		if err != nil {
			a.log.Error().Stack().Err(err).Msg("error happened while reading kafka messages")
		}
	}

	return nil
}
