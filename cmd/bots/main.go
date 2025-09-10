package main

import (
	"github.com/goriiin/kotyari-bots_backend/internal/apps/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"log"
)

var Config = "bots-config"

func main() {
	cfg, err := config.New[bots.BotsAppConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app := bots.NewApp(cfg)
	if err = app.Run(); err != nil {
		log.Fatalf("bots app exited with error: %v", err)
	}
}
