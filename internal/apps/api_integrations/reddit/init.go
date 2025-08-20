package reddit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery/api_integrations/reddit"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/producer"
	"github.com/goriiin/kotyari-bots_backend/internal/repo/api_integrations"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"github.com/segmentio/kafka-go"
)

type RedditAPIDelivery interface {
	Run() error
}

type RedditAPIApp struct {
	delivery  RedditAPIDelivery
	appConfig RedditAppConfig
}

func NewRedditAPIApp(config *RedditAppConfig) (*RedditAPIApp, error) {
	pgxPool, err := postgres.GetPool(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}

	kafkaProducer := producer.NewKafkaProducer(&config.Kafka)

	redditRepo := api_integrations.NewAPIIntegrationsRepo(pgxPool)
	redditDelivery := reddit.NewRedditApiIntegration(&http.Client{}, redditRepo,
		time.Duration(config.API.PollInterval)*time.Minute, kafkaProducer)

	return &RedditAPIApp{
		delivery:  redditDelivery,
		appConfig: *config,
	}, nil
}

func (r *RedditAPIApp) Run() error {
	return r.delivery.Run()
}

func TestRead() error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "reddit_messages",
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
	})

	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		fmt.Println(string(message.Value))
	}
}
