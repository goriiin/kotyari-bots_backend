package example

import (
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const local = "config-example"

func main() {
	cfg, _ := config.New[config.AppConfig]()

	fmt.Printf("окружение: %s\n", cfg.GetEnvironment())
	fmt.Printf("конфигурация: {%v}, dname: %s, dpass: %s. duser: %s", cfg.API, cfg.Database.Name, cfg.Database.Password, cfg.Database.User)

	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[config.AppConfig](local)
		if err != nil {
			return
		}

		cfg = newCfg
		fmt.Printf("новая: %s:%d\n", cfg.API.Host, cfg.API.Port)
	})
}
