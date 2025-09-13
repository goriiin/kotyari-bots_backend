package bots

import (
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type configAPI struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type BotsAppConfig struct {
	config.ConfigBase
	API             configAPI       `mapstructure:"bots_api" env:"BOTS"`
	Database        postgres.Config `mapstructure:"bots_database" env:"BOTS"`
	ProfilesSvcAddr string          `mapstructure:"profiles_svc_addr" env:"BOTS"`
}
