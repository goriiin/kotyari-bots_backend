package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-faster/errors"
	"golang.org/x/net/proxy"
)

const tcpNetwork = "tcp"

type Proxy struct {
	config *ProxyConfig
	dialer proxy.Dialer
}

func NewProxy(config *ProxyConfig) (*Proxy, error) {
	connUrl := fmt.Sprintf("%s:%v", config.ProxyAPI.Host, config.ProxyAPI.Port)

	dialer, err := proxy.SOCKS5(tcpNetwork, connUrl, nil, proxy.Direct)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proxy dialer")
	}

	return &Proxy{
		config: config,
		dialer: dialer,
	}, nil
}

func (p *Proxy) UseProxy(client *http.Client) *http.Client {
	client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return p.dialer.Dial(network, addr)
		},
	}

	return client
}
