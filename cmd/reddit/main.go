package main

import (
	"fmt"
	"log"

	redditapp "github.com/kotyari-bots_backend/internal/apps/api_integrations/reddit"
	"github.com/kotyari-bots_backend/pkg/config"
)

const local = "local-config"

func main() {

	cfg, _ := config.New[redditapp.RedditAppConfig]()

	fmt.Printf("окружение: %s\n", cfg.GetEnvironment())
	fmt.Println("конфигурация: ", cfg.API, cfg.Database, cfg.Kafka)

	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[redditapp.RedditAppConfig](local)
		if err != nil {
			return
		}

		cfg = newCfg

		// TODO: ??
		err = cfg.Validate()
		if err != nil {
			return
		}

		fmt.Printf("новая: %s:%d\n", cfg.API.Host, cfg.API.Port)
	})

	app, err := redditapp.NewRedditAPIApp(cfg)
	if err != nil {
		log.Fatalf("failed to init reddit app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("reddit app exited with error: %v", err)
	}
}
