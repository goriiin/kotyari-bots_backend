package aggregator

import (
	"context"

	aggregatorDelivery "github.com/goriiin/kotyari-bots_backend/internal/delivery/aggregator"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/aggregator"
	aggregatorService "github.com/goriiin/kotyari-bots_backend/internal/usecase/aggregator"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type AggregatorDelivery interface {
	Run() error
}

type AggregatorApp struct {
	delivery AggregatorDelivery
	config   AggregatorAppConfig
}

func NewAggregatorApp(config *AggregatorAppConfig) (*AggregatorApp, error) {
	pgxPool, err := postgres.GetPool(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}

	kafkaConsumer := consumer.NewKafkaConsumer(&config.Kafka)

	aggregatorRepo := aggregator.NewAggregatorRepo(pgxPool)
	aggregatorUseCase := aggregatorService.NewAggregatorService(aggregatorRepo)
	aggregatorDelivery := aggregatorDelivery.NewAggregatorDelivery(kafkaConsumer, aggregatorUseCase)

	return &AggregatorApp{
		delivery: aggregatorDelivery,
		config:   *config,
	}, nil
}

func (a *AggregatorApp) Run() error {
	return a.delivery.Run()
}
