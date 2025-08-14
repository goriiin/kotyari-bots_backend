package config

import "fmt"

type HTTPServerConfig struct {
	Host string `mapstructure:"host" env:"HOST"`
	Port int    `mapstructure:"port" env:"PORT"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" env:"HOST"`
	Port     int    `mapstructure:"port" env:"PORT"`
	Name     string `mapstructure:"name" env:"NAME"`
	User     string `mapstructure:"user" env:"USER"`
	Password string `mapstructure:"password" env:"PASSWORD"`
}

type AppConfig struct {
	ConfigBase
	API      HTTPServerConfig `mapstructure:"api" envPrefix:"API"`
	Database DatabaseConfig   `mapstructure:"database" envPrefix:"DATABASE"`
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
