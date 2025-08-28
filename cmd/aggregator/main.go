package main

import (
	"log"

	"github.com/goriiin/kotyari-bots_backend/internal/apps/aggregator"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

const local = "local-config"

func main() {
	cfg, _ := config.New[aggregator.AggregatorAppConfig]()
	config.WatchConfig(func() {
		newCfg, err := config.NewWithConfig[aggregator.AggregatorAppConfig](local)
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
		log.Fatalf("aggregator app exited with error: %v", err)
	}
}
