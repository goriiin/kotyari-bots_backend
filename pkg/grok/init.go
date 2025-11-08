package grok

import (
	"net/http"

	proxyPkg "github.com/goriiin/kotyari-bots_backend/pkg/proxy"
)

const defaultModel = "grok-3-mini"

type GrokClient struct {
	config *GrokClientConfig
	// log
	httpClient *http.Client
	model      string
}

type GrokClientOption func(*GrokClient)

func WithModel(model string) GrokClientOption {
	return func(c *GrokClient) {
		if model != "" {
			c.model = model
		}
	}
}

func NewGrokClient(config *GrokClientConfig, proxyCfg *proxyPkg.ProxyConfig, opts ...GrokClientOption) (*GrokClient, error) {
	proxy, err := proxyPkg.NewProxy(proxyCfg)
	if err != nil {
		return nil, err
	}
	httpClient := proxy.UseProxy(&http.Client{})
	httpClient.Timeout = config.Timeout

	client := &GrokClient{
		config:     config,
		httpClient: httpClient,
		model:      defaultModel,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}
