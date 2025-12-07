package otvet

import (
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const OtvetBaseURL = "https://otvet.mail.ru"

type OtvetClientConfig struct {
	config.ConfigBase
	AuthToken string        `mapstructure:"auth_token" env:"OTVET_AUTH_TOKEN"`
	Timeout   time.Duration `mapstructure:"request_timeout"`
}

func (o *OtvetClientConfig) Validate() error {
	if o.AuthToken == "" {
		return fmt.Errorf("missing auth token")
	}

	if o.Timeout == 0 {
		o.Timeout = 30 * time.Second
	}

	return nil
}

