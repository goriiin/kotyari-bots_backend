package kafka

import (
	"context"
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/segmentio/kafka-go"
)

type Command string

// Envelope TODO: Можно постараться сделать более общим, но пока сойдет
type Envelope struct {
	Command       Command `json:"command"`
	EntityID      string  `json:"entity_id"`
	Payload       []byte  `json:"payload"`
	CorrelationID string  `json:"correlation_id"`
	Attempt       int     `json:"attempt,omitempty"` // Пока пусть будет
}

type CommittableMessage struct {
	Msg   kafka.Message
	Ack   func(ctx context.Context) error
	Nack  func(ctx context.Context, err error) error
	Reply func(ctx context.Context, body []byte) error
}

type KafkaConfig struct {
	config.ConfigBase
	Kind    string   `mapstructure:"kind"`
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

func (k *KafkaConfig) Validate() error {
	if len(k.Brokers) == 0 {
		return fmt.Errorf("brokers cannot be empty")
	}

	if k.Topic == "" {
		return fmt.Errorf("topic should be presented in config")
	}

	if k.Kind == "consumer" && k.GroupID == "" {
		return fmt.Errorf("group ID should be presented in consumer config")
	}

	return nil
}

func EnsureTopicCreated(broker, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}
