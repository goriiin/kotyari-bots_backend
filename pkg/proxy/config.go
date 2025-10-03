package proxy

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type proxyAPIConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ProxyConfig struct {
	config.ConfigBase
	ProxyAPI proxyAPIConfig `mapstructure:"proxy_server"`
}

func (p *ProxyConfig) Validate() error {
	if p.ProxyAPI.Port < 1 || p.ProxyAPI.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", p.ProxyAPI.Port)
	}

	if p.ProxyAPI.Host == "" {
		return fmt.Errorf("proxy host should be presented in config")
	}

	return nil
}
