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

			ctx, cancel := context.WithTimeout(ctx, batchTimeout)
			defer cancel()

			for len(messages) < batchSize {
				message, err := k.reader.ReadMessage(ctx)
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						break
					}

					if errors.Is(err, context.Canceled) {
						k.log.Warn("kafka is shutting down", err)
						return
					}

					k.log.Error(err, false, "unexpected error happened")
					return
				}

				messages = append(messages, message)
			}

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
