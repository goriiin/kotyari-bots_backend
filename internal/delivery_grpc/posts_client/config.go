package posts_client

import (
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type postsGRPCClientApiCfg struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type PostsGRPCClientAppConfig struct {
	config.ConfigBase
	API          postsGRPCClientApiCfg `mapstructure:"posts_grpc_api"`
	BotsAddr     string                `mapstructure:"bots_grpc_addr"`
	ProfilesAddr string                `mapstructure:"bots_grpc_addr"`
	Timeout      time.Duration         `mapstructure:"dial_timeout"`
}

func (p *PostsGRPCClientAppConfig) Validate() error {
	if p.API.Port < 1 || p.API.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", p.API.Port)
	}

	if p.API.Host == "" {
		return fmt.Errorf("host should be presented in config")
	}

	if p.BotsAddr == "" || p.ProfilesAddr == "" {
		return fmt.Errorf("bots and profile servers adresses should be presented in config")
	}

	if p.Timeout == 0 {
		p.Timeout = time.Duration(5) * time.Second
	}

	return nil
}
