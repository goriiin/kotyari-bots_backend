package consumer

import (
	"context"
	"errors"
	"time"

	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

const (
	batchSize    = 20
	batchTimeout = 20 * time.Second
)

type KafkaConsumer struct {
	// TODO: logs
	reader *kafka.Reader
}

func NewKafkaConsumer(config *kafkaConfig.KafkaConfig) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: config.Brokers,
			Topic:   config.Topic,
			GroupID: config.GroupID,
		}),
	}
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
						// TODO: logs - shutdown

						return
					}

					// TODO: log err

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

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return batches
}
