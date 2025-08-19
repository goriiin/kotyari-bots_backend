package postgres

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/pkg/utils"
)

func ToUrl(conf Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
		utils.GetValueOrDefault(conf.SSLMode, defaultSSLMode),
	)
}
