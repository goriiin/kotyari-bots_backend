package config

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type HTTPServerConfig struct {
	Host string `mapstructure:"host" env:"HOST"`
	Port int    `mapstructure:"port" env:"PORT"`
}

// AppConfig EXAMPLE
type AppConfig struct {
	ConfigBase
	API      HTTPServerConfig `mapstructure:"api" envPrefix:"API"`
	Database postgres.Config  `mapstructure:"database" envPrefix:"DATABASE"`
}

func (c *AppConfig) Validate() error {
	if c.API.Port < 1 || c.API.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", c.API.Port)
	}

	if c.IsProduction() && c.Database.Password == "" {
		return fmt.Errorf("database password is required in production")
	}

	return nil
}
