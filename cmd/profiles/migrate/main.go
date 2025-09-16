package main

import (
	"github.com/goriiin/kotyari-bots_backend/internal/apps/profiles"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/migrate"
	"log"
)

func main() {
	cfg, err := config.New[profiles.ProfilesAppConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	migrate.Run(cfg.Database)
}
