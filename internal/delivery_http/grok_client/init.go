package grok_client

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/pkg/utils"
)

type GrokClient struct {
	config *GrokClientConfig
	//log
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

func (c *GrokClient) TestRequest() {
	targetURL := "https://ipinfo.io/ip"
	log.Printf("Making test request to %s via VLESS proxy...", targetURL)

	resp, err := c.httpClient.Get(targetURL)
	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	log.Printf("Request successful! Status: %s", resp.Status)
	log.Printf("Proxy IP Address: %s", string(body))
}
