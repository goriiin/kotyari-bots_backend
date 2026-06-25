package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/posts_command_producer"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

func main() {
	cfg, err := config.New[posts_command_producer.PostsCommandProducerConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := posts_command_producer.NewPostsCommandProducerApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	runErr := make(chan error, 1)
	go func() {
		runErr <- app.Run()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-runErr:
		if err != nil {
			log.Println("run error:", err)
		}
	case s := <-sig:
		log.Printf("received signal %s, shutting down", s)
	}

	if err := app.Close(); err != nil {
		// Log, don't Fatal: os.Exit here would skip any remaining cleanup.
		log.Println("close error:", err)
	}
}
