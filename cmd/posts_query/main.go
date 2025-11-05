package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

func main() {
	cfg, err := config.New[posts_query.PostsQueryConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := posts_query.NewPostsQueryApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
