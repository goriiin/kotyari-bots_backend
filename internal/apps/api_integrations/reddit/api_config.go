package reddit

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type redditAPIConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	PollInterval int    `mapstructure:"poll_interval"`
}
type RedditAppConfig struct {
	config.ConfigBase
	API      redditAPIConfig   `mapstructure:"reddit_api"`
	Database postgres.Config   `mapstructure:"reddit_database"`
	Kafka    kafka.KafkaConfig `mapstructure:"reddit_producer"`
}

func (c *RedditAppConfig) Validate() error {
	if err := c.Kafka.Validate(); err != nil {
		return err
	}

	if c.API.Port < 1 || c.API.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", c.API.Port)
	}
	if c.IsProduction() && c.Database.Password == "" {
		return fmt.Errorf("database password is required in production")
	}

	if c.API.PollInterval <= 0 {
		return fmt.Errorf("interval with value %d is not allowed", c.API.PollInterval)
	}

	return nil
}
