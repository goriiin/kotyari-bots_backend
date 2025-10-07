package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/proxy"
)

const local = "docker-config"

func main() {
	postsCfg, _ := config.New[posts.PostsAppCfg]()

	config.WatchConfig(func() {
		newPostsCfg, err := config.NewWithConfig[posts.PostsAppCfg](local)
		if err != nil {
			log.Fatalf("error parsing posts config in runtime: %s", err.Error())
			return
		}

		postsCfg = newPostsCfg
	})

	proxyCfg, _ := config.New[proxy.ProxyConfig]()
	config.WatchConfig(func() {
		newProxyCfg, err := config.NewWithConfig[proxy.ProxyConfig](local)
		if err != nil {
			log.Fatalf("error parsing proxy config in runtime: %s", err.Error())
			return
		}

		proxyCfg = newProxyCfg
	})

	posts, err := posts.NewPostsApp(postsCfg, proxyCfg)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	if err := posts.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
