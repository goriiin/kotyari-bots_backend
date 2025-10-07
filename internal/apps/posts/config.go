package posts

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_client"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type PostsApiCfg struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type PostsAppCfg struct {
	config.ConfigBase
	API        PostsApiCfg                           `mapstructure:"posts_api"`
	GrpcClient posts_client.PostsGRPCClientAppConfig `mapstructure:"posts_grpc"`
	Database   postgres.Config                       `mapstructure:"posts_db"`
	GrokCfg    grok.GrokClientConfig                 `mapstructure:"posts_grok"`
}

func (p *PostsAppCfg) Validate() error {
	if err := p.GrpcClient.Validate(); err != nil {
		return err
	}
	if err := p.GrokCfg.Validate(); err != nil {
		return err
	}

	if p.API.Port < 1 || p.API.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", p.API.Port)
	}

	if p.API.Host == "" {
		return fmt.Errorf("host should be presented in config")
	}

	return nil
}
