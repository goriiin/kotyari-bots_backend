package posts_command_consumer

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_command_consumer"
)

func main() {
	app, err := posts_command_consumer.NewPostsCommandConsumer()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
