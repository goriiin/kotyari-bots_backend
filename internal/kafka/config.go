package kafka

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

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
