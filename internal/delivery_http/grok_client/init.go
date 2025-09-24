package grok_client

import (
	"context"
	"net"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/pkg/utils"
)

type GrokClient struct {
	config *GrokClientConfig
	// log
	httpClient *http.Client
	xray       *utils.XrayCoreInstance
}

func NewGrokClient(config *GrokClientConfig) (*GrokClient, error) {
	vlessParams, err := utils.ParseVlessConfig(config.ProxyURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse vless url")
	}

	xrayInstance, err := utils.NewXrayCoreInstance(vlessParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create xray instance")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return xrayInstance.Dialer.Dial(network, addr)
			},
		},
		Timeout: config.Timeout,
	}

	return &GrokClient{
		config:     config,
		httpClient: httpClient,
		xray:       xrayInstance,
	}, nil
}

func (c *GrokClient) Close() error {
	if c.xray != nil {
		return c.xray.Instance.Close()
	}

	return nil
}
