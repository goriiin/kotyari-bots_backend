package aggregator

import (
	"context"
	"fmt"
)

func (a *AggregatorDelivery) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("here")
	messageBatches := a.consumer.ReadBatches(ctx)

	for batch := range messageBatches {
		err := a.manager.AddTopics(ctx, batch)
		if err != nil {
			// TODO: log
			return err
		}
	}

	return nil
}
