package posts_query

import (
	"github.com/goriiin/kotyari-bots_backend/internal/adapters/auth"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type configAPI struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type PostsQueryConfig struct {
	config.ConfigBase
	API      configAPI       `mapstructure:"posts_query_api"`
	Database postgres.Config `mapstructure:"posts_database"`
	Auth     auth.Config     `mapstructure:"auth_grpc"`
}
