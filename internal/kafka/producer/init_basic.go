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
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (k *KafkaProducer) Publish(ctx context.Context, message kafka.Message) error {
	return k.writer.WriteMessages(ctx, message)
}

func (k *KafkaProducer) Close() error {
	return k.writer.Close()
}
