package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/profiles"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

var Config = "profiles-config"

func main() {
	cfg, err := config.New[profiles.ProfilesAppConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app := profiles.NewApp(cfg)
	if err = app.Run(); err != nil {
		log.Fatalf("bots app exited with error: %v", err)
	}
}
