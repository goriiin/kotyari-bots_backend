package utils

import (
	"encoding/json"
	"fmt"

	"github.com/go-faster/errors"
	xnet "github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	"golang.org/x/net/proxy"

	_ "github.com/xtls/xray-core/main/distro/all"
)

// Не уверен, нужно ли это (как минимум порт и адрес) выносить в конфиг, или оставить константами
const (
	localSocksPort    = 8000
	localSocksAddress = "127.0.0.1"
	realitySec        = "reality"
	defaultTag        = "proxy"
	vlessProtocol     = "vless"
)

type XrayCoreInstance struct {
	Instance *core.Instance
	Dialer   proxy.Dialer
}

func NewXrayCoreInstance(vlessParams *VlessConfig) (*XrayCoreInstance, error) {
	userBytes, err := json.Marshal(newVlessUserSetting(vlessParams))
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal vless user settings")
	}

	vlessOutboundSettings := &conf.VLessOutboundConfig{
		Vnext: []*conf.VLessOutboundVnext{
			{
				Address: &conf.Address{
					Address: xnet.ParseAddress(vlessParams.Address),
				},
				Port:  vlessParams.Port,
				Users: []json.RawMessage{userBytes},
			},
		},
	}

	rawVlessSettings, err := json.Marshal(vlessOutboundSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal vless outbound settings")
	}

	jsonRawVlessSettings := json.RawMessage(rawVlessSettings)
	networkProtocol := conf.TransportProtocol(vlessParams.Network)
	streamConfig := &conf.StreamConfig{
		Network:  &networkProtocol,
		Security: vlessParams.Security,
	}
	if vlessParams.Security == realitySec {
		streamConfig.REALITYSettings = &conf.REALITYConfig{
			ServerName:  vlessParams.SNI,
			Fingerprint: vlessParams.Fingerprint,
			ShortId:     vlessParams.ShortID,
			PublicKey:   vlessParams.PublicKey,
		}
	}

	socksInboundSettings := &conf.SocksServerConfig{
		AuthMethod: "noauth",
		UDP:        true,
	}
	rawSocksSettings, err := json.Marshal(socksInboundSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal socks inbound settings")
	}

	jsonRawSocksSettings := json.RawMessage(rawSocksSettings)
	xrayConfig := &conf.Config{
		InboundConfigs: []conf.InboundDetourConfig{
			{
				Tag:      "socks-in",
				Protocol: "socks",
				ListenOn: &conf.Address{
					Address: xnet.ParseAddress(localSocksAddress),
				},
				PortList: &conf.PortList{
					Range: []conf.PortRange{
						{From: localSocksPort, To: localSocksPort},
					}},
				Settings: &jsonRawSocksSettings,
			},
		},
		OutboundConfigs: []conf.OutboundDetourConfig{
			{
				Tag:           defaultTag,
				Protocol:      vlessProtocol,
				Settings:      &jsonRawVlessSettings,
				StreamSetting: streamConfig,
			},
		},
	}

	coreConfig, err := xrayConfig.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build xray config")
	}

	xrayInstance, err := core.New(coreConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create xray instance")
	}
	if err := xrayInstance.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start xray instance")
	}

	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", localSocksAddress, localSocksPort), nil, proxy.Direct)
	if err != nil {
		xrayInstance.Close()
		return nil, errors.Wrap(err, "failed to create socks5 dialer")
	}

	return &XrayCoreInstance{
		Instance: xrayInstance,
		Dialer:   dialer,
	}, nil
}
