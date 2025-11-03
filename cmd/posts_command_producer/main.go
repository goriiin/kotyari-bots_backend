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
