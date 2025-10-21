package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaRequestReplyProducer struct {
	writer     *kafka.Writer
	config     *kafkaConfig.KafkaConfig
	replyTopic string
	replyGroup string
}

func NewKafkaRequestReplyProducer(config *kafkaConfig.KafkaConfig, replyTopic, replyGroup string) *KafkaRequestReplyProducer {
	return &KafkaRequestReplyProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(config.Brokers...),
			Topic:                  config.Topic,
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
		},
		replyTopic: replyTopic,
		replyGroup: replyGroup,
	}
}

func (p *KafkaRequestReplyProducer) Publish(ctx context.Context, env kafkaConfig.Envelope) error {
	env.CorrelationID = uuid.NewString()
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(env.EntityID),
		Value: b,
		Headers: []kafka.Header{
			{Key: "correlation_id", Value: []byte(env.CorrelationID)},
			{Key: "command", Value: []byte(env.Command)},
			{Key: "reply_to", Value: []byte(p.replyTopic)},
		},
	})
}

func (p *KafkaRequestReplyProducer) Request(ctx context.Context, env kafkaConfig.Envelope, timeout time.Duration) ([]byte, error) {
	if err := p.Publish(ctx, env); err != nil {
		return nil, err
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  p.config.Brokers,
		Topic:    p.replyTopic,
		GroupID:  p.replyGroup,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return nil, err
		}
		if getHeader(m, "correlation_id") == env.CorrelationID {
			return m.Value, nil
		}
	}
}

func getHeader(m kafka.Message, key string) string {
	for _, h := range m.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}
