package main

import (
	"fmt"
	"github.com/kotyari-bots_backend/pkg/config"
)

const local = "config-local"

func main() {
	cfg, _ := config.New[config.AppConfig]()

	fmt.Printf("окружение: %s\n", cfg.GetEnvironment())
	fmt.Println("конфигурация: ", cfg.API, cfg.Database)

	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[config.AppConfig](local)
		if err != nil {
			return
		}

		cfg = newCfg
		fmt.Printf("новая: %s:%d\n", cfg.API.Host, cfg.API.Port)
	})

	fmt.Scan()
}
