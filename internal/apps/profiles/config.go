package profiles

import (
	"github.com/goriiin/kotyari-bots_backend/internal/adapters/auth"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type configAPI struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	GRPCPort int    `mapstructure:"grpc_port"`
}

type ProfilesAppConfig struct {
	config.ConfigBase
	API      configAPI       `mapstructure:"profiles_api"`
	Database postgres.Config `mapstructure:"profiles_database"`
	Auth     auth.Config     `mapstructure:"auth_grpc"`
}
