package grok

import (
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const GrokTargetUrl = "https://api.x.ai/v1/chat/completions"

type GrokClientConfig struct {
	config.ConfigBase
	ApiKey  string        `mapstructure:"api_key"`
	Timeout time.Duration `mapstructure:"request_timeout"`
}

func (g *GrokClientConfig) Validate() error {

	if g.ApiKey == "" {
		return fmt.Errorf("missing API key")
	}

	if g.Timeout == 0 {
		g.Timeout = 30 * time.Second
	}

	return nil
}
