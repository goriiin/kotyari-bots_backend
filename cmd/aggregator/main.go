package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/aggregator"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

func main() {
	cfg, err := config.New[aggregator.AggregatorAppConfig]()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := aggregator.NewAggregatorApp(cfg)
	if err != nil {
		log.Fatalf("failed to init aggregator app: %v", err)
	}

	if err := app.Run(); err != nil {
		app.Log.Fatal(err, true, "aggregator app exited with error")
	}
}
