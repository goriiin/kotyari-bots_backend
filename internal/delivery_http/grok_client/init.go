package grok_client

import (
	"net/http"

	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
	proxyPkg "github.com/goriiin/kotyari-bots_backend/pkg/proxy"
)

type GrokClient struct {
	config *grok.GrokClientConfig
	// log
	httpClient *http.Client
}

func NewGrokClient(config *grok.GrokClientConfig, proxyCfg *proxyPkg.ProxyConfig) (*GrokClient, error) {
	proxy, err := proxyPkg.NewProxy(proxyCfg)
	if err != nil {
		return nil, err
	}
	httpClient := proxy.UseProxy(&http.Client{})
	httpClient.Timeout = config.Timeout

	return &GrokClient{
		config:     config,
		httpClient: httpClient,
	}, nil
}
