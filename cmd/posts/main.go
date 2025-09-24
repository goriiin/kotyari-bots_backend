package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"

	_ "github.com/xtls/xray-core/proxy/socks"
	_ "github.com/xtls/xray-core/proxy/vless"
	_ "github.com/xtls/xray-core/proxy/vless/inbound"
	_ "github.com/xtls/xray-core/proxy/vless/outbound"
	_ "github.com/xtls/xray-core/transport/internet/reality"
	_ "github.com/xtls/xray-core/transport/internet/tcp"
	_ "github.com/xtls/xray-core/transport/internet/tls"
)

const local = "local-config"

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
