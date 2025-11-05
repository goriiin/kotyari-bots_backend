package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_command_producer"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

func main() {
	cfg, err := config.New[posts_command_producer.PostsCommandProducerConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := posts_command_producer.NewPostsCommandProducerApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func(app *posts_command_producer.PostsCommandProducerApp) {
		err := app.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(app)

	if err = app.Run(); err != nil {
		log.Println(err)
	}
}
