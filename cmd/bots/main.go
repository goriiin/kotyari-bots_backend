package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

var Config = "bots-local"

func main() {
	cfg, err := config.New[bots.AppConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app := bots.NewApp(cfg)
	if err = app.Run(); err != nil {
		log.Fatalf("bots app exited with error: %v", err)
	}
}
