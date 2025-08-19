package kafka

import (
	"fmt"

	"github.com/kotyari-bots_backend/pkg/config"
)

type KafkaConfig struct {
	config.ConfigBase
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

func (k *KafkaConfig) Validate() error {
	if len(k.Brokers) == 0 {
		return fmt.Errorf("brokers cannot be empty")
	}

	if k.Topic == "" {
		return fmt.Errorf("topic should be presented in config")
	}

	return nil
}
