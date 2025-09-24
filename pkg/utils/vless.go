package utils

import (
	"net/url"
	"strconv"

	"github.com/go-faster/errors"
)

type VlessConfig struct {
	UserID      string
	Address     string
	Port        uint16
	Encryption  string
	Flow        string
	Network     string
	Security    string
	SNI         string
	Fingerprint string
	PublicKey   string
	ShortID     string
}

const (
	defaultVlessFlow = "xtls-rprx-vision"
)

type vlessUserSetting struct {
	ID         string `json:"id"`
	Flow       string `json:"flow"`
	Encryption string `json:"encryption"`
}

func newVlessUserSetting(vlessCfg *VlessConfig) *vlessUserSetting {
	return &vlessUserSetting{
		ID:         vlessCfg.UserID,
		Flow:       defaultVlessFlow,
		Encryption: vlessCfg.Encryption,
	}
}

func ParseVlessConfig(vlessURL string) (*VlessConfig, error) {
	u, err := url.Parse(vlessURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid VLESS URL: %w")
	}

	if u.Scheme != "vless" {
		return nil, errors.Errorf("not a vless scheme: %s", u.Scheme)
	}

	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid port: %d", port)
	}

	q := u.Query()

	return &VlessConfig{
		UserID:      u.User.Username(),
		Address:     u.Hostname(),
		Port:        uint16(port),
		Encryption:  q.Get("encryption"),
		Flow:        q.Get("flow"),
		Network:     q.Get("type"),
		Security:    q.Get("security"),
		SNI:         q.Get("sni"),
		Fingerprint: q.Get("fp"),
		PublicKey:   q.Get("pbk"),
		ShortID:     q.Get("sid"),
	}, nil
}
