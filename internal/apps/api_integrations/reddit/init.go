package reddit

import (
	"context"
	"net/http"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/api_integrations/reddit"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/producer"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/api_integrations"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

const serviceName = "reddit-app"

type RedditAPIDelivery interface {
	Run() error
}

type RedditAPIApp struct {
	delivery  RedditAPIDelivery
	appConfig RedditAppConfig
	Log       logger.Logger
}

func NewRedditAPIApp(config *RedditAppConfig) (*RedditAPIApp, error) {
	pgxPool, err := postgres.GetPool(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger(serviceName, &config.ConfigBase)

	kafkaProducer := producer.NewKafkaProducer(&config.Kafka)

	redditRepo := api_integrations.NewAPIIntegrationsRepo(*log, pgxPool)
	redditDelivery := reddit.NewRedditApiIntegration(&http.Client{}, redditRepo,
		time.Duration(config.API.PollInterval)*time.Minute, kafkaProducer, *log)

	return &RedditAPIApp{
		delivery:  redditDelivery,
		appConfig: *config,
		Log:       *log,
	}, nil
}

func (r *RedditAPIApp) Run() error {
	return r.delivery.Run()
}
