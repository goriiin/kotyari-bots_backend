package consumer

import (
	"context"
	"errors"
	"time"

	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/segmentio/kafka-go"
)

const (
	batchSize    = 20
	batchTimeout = 20 * time.Second
)

type KafkaConsumer struct {
	log    *logger.Logger
	reader *kafka.Reader
}

func NewKafkaConsumer(log *logger.Logger, config *kafkaConfig.KafkaConfig) *KafkaConsumer {
	return &KafkaConsumer{
		log: log,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: config.Brokers,
			Topic:   config.Topic,
			GroupID: config.GroupID,
		}),
	}
}

func (k *KafkaConsumer) GetMessage(ctx context.Context) (kafka.Message, error) {
	return k.reader.ReadMessage(ctx)
}

func (k *KafkaConsumer) ReadBatches(ctx context.Context) <-chan []kafka.Message {
	batches := make(chan []kafka.Message)

	go func() {
		defer close(batches)

		for {
			var messages []kafka.Message

			// Per-batch timeout context. Cancel it at the end of each iteration
			// instead of deferring (a defer here would accumulate one live timer
			// per loop iteration for the lifetime of this goroutine).
			batchCtx, cancel := context.WithTimeout(ctx, batchTimeout)

			for len(messages) < batchSize {
				message, err := k.reader.ReadMessage(batchCtx)
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						break
					}

					if errors.Is(err, context.Canceled) {
						k.log.Warn("kafka is shutting down", err)
						cancel()
						return
					}

					k.log.Error(err, false, "unexpected error happened")
					cancel()
					return
				}

				messages = append(messages, message)
			}

			cancel()

			if len(messages) > 0 {
				select {
				case batches <- messages:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return batches
}
