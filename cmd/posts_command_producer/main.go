package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_command_producer"
)

func main() {
	app, err := posts_command_producer.NewPostsCommandProducerApp()
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
