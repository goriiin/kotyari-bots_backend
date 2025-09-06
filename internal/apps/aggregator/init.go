package aggregator

import (
	"context"

	aggregatorDelivery "github.com/goriiin/kotyari-bots_backend/internal/delivery_http/aggregator"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/aggregator"
	aggregatorService "github.com/goriiin/kotyari-bots_backend/internal/usecase/aggregator"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

const serviceName = "aggregator-app"

type AggregatorDelivery interface {
	Run() error
}

type AggregatorApp struct {
	delivery AggregatorDelivery
	config   AggregatorAppConfig
	Log      *logger.Logger
}

func NewAggregatorApp(config *AggregatorAppConfig) (*AggregatorApp, error) {
	pgxPool, err := postgres.GetPool(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger(serviceName, &config.ConfigBase)

	kafkaConsumer := consumer.NewKafkaConsumer(log, &config.Kafka)

	aggregatorRepo := aggregator.NewAggregatorRepo(log, pgxPool)
	aggregatorUseCase := aggregatorService.NewAggregatorService(log, aggregatorRepo)
	aggregatorDel := aggregatorDelivery.NewAggregatorDelivery(log, kafkaConsumer, aggregatorUseCase)

	return &AggregatorApp{
		Log:      log,
		delivery: aggregatorDel,
		config:   *config,
	}, nil
}

func (a *AggregatorApp) Run() error {
	return a.delivery.Run()
}
