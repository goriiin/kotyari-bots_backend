package aggregator

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/segmentio/kafka-go"
)

type AggregatorUseCase interface {
	AddTopics(ctx context.Context, topics []kafka.Message) error
}

type MessageConsumer interface {
	ReadBatches(ctx context.Context) <-chan []kafka.Message
}

type AggregatorDelivery struct {
	log      *logger.Logger
	consumer MessageConsumer
	manager  AggregatorUseCase
}

func NewAggregatorDelivery(log *logger.Logger, consumer MessageConsumer, manager AggregatorUseCase) *AggregatorDelivery {
	return &AggregatorDelivery{
		log:      log,
		consumer: consumer,
		manager:  manager,
	}
}
