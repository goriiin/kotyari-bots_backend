package grok_client

import (
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type grokAPIConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GrokClientConfig struct {
	config.ConfigBase
	API      grokAPIConfig `mapstructure:"grok_server"`
	ApiKey   string        `mapstructure:"api_key"`
	ProxyURL string        `mapstructure:"proxy_url"`
	Timeout  time.Duration `mapstructure:"request_timeout"`
}

func (g *GrokClientConfig) Validate() error {
	if g.API.Port < 1 || g.API.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", g.API.Port)
	}

	if g.API.Host == "" {
		return fmt.Errorf("host should be presented in config")
	}

	if g.ApiKey == "" {
		return fmt.Errorf("missing API key")
	}

	if g.ProxyURL == "" {
		return fmt.Errorf("missing proxy url")
	}

	if g.Timeout == 0 {
		g.Timeout = 30 * time.Second
	}

	return nil
}
