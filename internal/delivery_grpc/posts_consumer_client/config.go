package posts_consumer_client

import (
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type PostsConsGRPCClientConfig struct {
	config.ConfigBase
	PostsAddr string        `mapstructure:"posts_addr"`
	Timeout   time.Duration `mapstructure:"dial_timeout"`
}

func (p *PostsConsGRPCClientConfig) Validate() error {
	if p.PostsAddr == "" {
		return fmt.Errorf("bots, profile and posts servers adresses should be presented in config")
	}

	if p.Timeout == 0 {
		p.Timeout = time.Duration(5) * time.Second
	}

	return nil
}
