package consumer

import (
	kafkaConfig "github.com/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	// TODO: logs
	reader *kafka.Reader
}

func NewKafkaConsumer(config *kafkaConfig.KafkaConfig) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Brokers,
		Topic:   config.Topic,
	})

	return &KafkaConsumer{
		reader: reader,
	}
}
