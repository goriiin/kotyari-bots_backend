package producer

import (
	"context"

	kafka2 "github.com/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	// TODO: logs
	writer *kafka.Writer
}

func NewKafkaProducer(config *kafka2.KafkaConfig) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  config.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{
		writer: writer,
	}
}

func (k *KafkaProducer) Publish(ctx context.Context, message []byte) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}
