package proxy

import (
	"fmt"

	"github.com/go-faster/errors"
	"golang.org/x/net/proxy"
)

const tcpNetwork = "tcp"

type Proxy struct {
	config *ProxyConfig
	Dialer proxy.Dialer
}

func NewProxy(config *ProxyConfig) (*Proxy, error) {
	connUrl := fmt.Sprintf("%s:%v", config.ProxyAPI.Host, config.ProxyAPI.Port)

	dialer, err := proxy.SOCKS5(tcpNetwork, connUrl, nil, proxy.Direct)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proxy dialer")
	}

	return &Proxy{
		config: config,
		Dialer: dialer,
	}, nil
}
