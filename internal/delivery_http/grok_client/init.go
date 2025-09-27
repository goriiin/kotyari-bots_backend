package grok_client

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
	"golang.org/x/net/proxy"
)

type GrokClient struct {
	config *grok.GrokClientConfig
	// log
	httpClient *http.Client
}

func NewGrokClient(config *grok.GrokClientConfig) (*GrokClient, error) {
	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", config.ProxyAPI.Host, config.ProxyAPI.Port), nil, proxy.Direct)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dialer")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
		Timeout: config.Timeout,
	}

	return &GrokClient{
		config:     config,
		httpClient: httpClient,
	}, nil
}
