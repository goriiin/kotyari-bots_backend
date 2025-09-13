package aggregator

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type AggregatorUseCase interface {
	AddTopics(ctx context.Context, topics []kafka.Message) error
}

type MessageConsumer interface {
	ReadBatches(ctx context.Context) <-chan []kafka.Message
}

type AggregatorDelivery struct {
	consumer MessageConsumer
	manager  AggregatorUseCase
}

func NewAggregatorDelivery(consumer MessageConsumer, manager AggregatorUseCase) *AggregatorDelivery {
	return &AggregatorDelivery{
		consumer: consumer,
		manager:  manager,
	}
}
