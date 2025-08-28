package main

import (
	"log"

	redditapp "github.com/goriiin/kotyari-bots_backend/internal/apps/api_integrations/reddit"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const local = "local-config"

func main() {
	cfg, _ := config.New[redditapp.RedditAppConfig]()

	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[redditapp.RedditAppConfig](local)
		if err != nil {
			log.Fatalf("error parsing config in runtime: %s", err.Error())
			return
		}

		cfg = newCfg
	})

	app, err := redditapp.NewRedditAPIApp(cfg)
	if err != nil {
		log.Fatalf("failed to init reddit app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("reddit app exited with error: %v", err)
	}
}
