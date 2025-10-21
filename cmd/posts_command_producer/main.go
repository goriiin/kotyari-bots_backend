package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_command_producer"
)

func main() {
	app := posts_command_producer.NewPostsCommandProducerApp()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
