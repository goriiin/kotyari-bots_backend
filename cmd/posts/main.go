package main

import (
	"log"

	postsclient "github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_client"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const local = "config-local"

func main() {

	cfg, _ := config.New[postsclient.PostsGRPCClientAppConfig]()

	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[postsclient.PostsGRPCClientAppConfig](local)
		if err != nil {
			log.Fatalf("error parsing config in runtime: %s", err.Error())
			return
		}

		cfg = newCfg
	})

	grpcClient, err := postsclient.NewPostsGRPCClient(cfg)
	if err != nil {
		log.Fatalf("failed to init grpc client grpcClient grpcClient: %v", err)
	}

	// TODO: move to defer func() later
	if err := grpcClient.Close(); err != nil {
		log.Printf("failed to close grpc client connections: %v", err)
	}
}
