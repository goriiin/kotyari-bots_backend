package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const local = "docker-config"

func main() {
	cfg, _ := config.New[posts.PostsAppCfg]()

	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[posts.PostsAppCfg](local)
		if err != nil {
			log.Fatalf("error parsing config in runtime: %s", err.Error())
			return
		}

		cfg = newCfg
	})

	_, err := posts.NewPostsApp(cfg)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
