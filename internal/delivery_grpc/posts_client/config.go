package posts_client

import (
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type PostsGRPCClientAppConfig struct {
	config.ConfigBase
	BotsAddr     string        `mapstructure:"bots_addr"`
	ProfilesAddr string        `mapstructure:"profiles_addr"`
	Timeout      time.Duration `mapstructure:"dial_timeout"`
}

func (p *PostsGRPCClientAppConfig) Validate() error {

	if p.BotsAddr == "" || p.ProfilesAddr == "" {
		return fmt.Errorf("bots and profile servers adresses should be presented in config")
	}

	if p.Timeout == 0 {
		p.Timeout = time.Duration(5) * time.Second
	}

	return nil
}
