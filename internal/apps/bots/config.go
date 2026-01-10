package bots

import (
	"github.com/goriiin/kotyari-bots_backend/internal/adapters/auth"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type configAPI struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type configGRPC struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type AppConfig struct {
	config.ConfigBase
	API             configAPI       `mapstructure:"bots_api" env:"BOTS"`
	GRPC            configGRPC      `mapstructure:"bots_grpc" env:"BOTS_GRPC"`
	Database        postgres.Config `mapstructure:"bots_database" env:"BOTS"`
	ProfilesSvcAddr string          `mapstructure:"profiles_svc_addr" env:"BOTS"`
	Auth            auth.Config     `mapstructure:"auth_grpc"`
}
