package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_command_consumer"
)

func main() {
	app, err := posts_command_consumer.NewPostsCommandConsumer()
	if err != nil {
		log.Fatal(err)
	}

	defer func(app *posts_command_consumer.PostsCommandConsumer) {
		err := app.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(app)

	if err = app.Run(); err != nil {
		log.Println(err)
	}
}
