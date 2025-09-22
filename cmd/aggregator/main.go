package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/aggregator"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

// const local = "local-config"
const docker = "docker-config"

func main() {
	cfg, _ := config.New[aggregator.AggregatorAppConfig]()
	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[aggregator.AggregatorAppConfig](docker)
		if err != nil {
			return
		}

		cfg = newCfg
	})

	app, err := aggregator.NewAggregatorApp(cfg)
	if err != nil {
		log.Fatalf("failed to init aggregator app: %v", err)
	}

	if err := app.Run(); err != nil {
		app.Log.Fatal(err, true, "aggregator app exited with error")
	}
}
