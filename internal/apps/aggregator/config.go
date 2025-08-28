package aggregator

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type aggregatorConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type AggregatorAppConfig struct {
	config.ConfigBase
	API      aggregatorConfig  `mapstructure:"aggregator_api"`
	Database postgres.Config   `mapstructure:"aggregator_database"`
	Kafka    kafka.KafkaConfig `mapstructure:"aggregator_consumer"`
}

func (a *AggregatorAppConfig) Validate() error {
	if err := a.Kafka.Validate(); err != nil {
		return err
	}

	if a.API.Port < 1 || a.API.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", a.API.Port)
	}
	if a.IsProduction() && a.Database.Password == "" {
		return fmt.Errorf("database password is required in production")
	}

	return nil
}
