package main

import (
	"log"

	redditapp "github.com/goriiin/kotyari-bots_backend/internal/apps/api_integrations/reddit"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

func main() {
	cfg, err := config.New[redditapp.RedditAppConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := redditapp.NewRedditAPIApp(cfg)
	if err != nil {
		log.Fatalf("failed to init reddit app: %v", err)
	}

	if err := app.Run(); err != nil {
		app.Log.Fatal(err, true, "reddit app exited with error")
	}
}
