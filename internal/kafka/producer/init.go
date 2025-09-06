package producer

import (
	"context"

	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(config *kafkaConfig.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(config.Brokers...),
			Topic:                  config.Topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (k *KafkaProducer) Publish(ctx context.Context, message []byte) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}
