package grok

import (
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const GrokTargetUrl = "https://api.x.ai/v1/chat/completions"

// TODO: Мне кажется, лучше выделить в отдельный базовый конфиг какой-то (что-то в духе config.BaseConfig)
// Потому что много где юзается эта структура с хостом и портом
type proxyAPIConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GrokClientConfig struct {
	config.ConfigBase
	ApiKey   string         `mapstructure:"api_key"`
	ProxyAPI proxyAPIConfig `mapstructure:"proxy_server"`
	Timeout  time.Duration  `mapstructure:"request_timeout"`
}

func (g *GrokClientConfig) Validate() error {
	if g.ProxyAPI.Port < 1 || g.ProxyAPI.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", g.ProxyAPI.Port)
	}

	if g.ProxyAPI.Host == "" {
		return fmt.Errorf("host should be presented in config")
	}

	if g.ApiKey == "" {
		return fmt.Errorf("missing API key")
	}

	if g.Timeout == 0 {
		g.Timeout = 30 * time.Second
	}

	return nil
}
