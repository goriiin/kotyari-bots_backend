package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_command_consumer"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/migrate"
)

func main() {
	cfg, err := config.New[posts_command_consumer.PostsCommandConsumerConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	migrate.Run(cfg.Database)
}
