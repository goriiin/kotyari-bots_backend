package reddit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kotyari-bots_backend/internal/delivery/api_integrations/reddit"
	"github.com/kotyari-bots_backend/internal/kafka/producer"
	"github.com/kotyari-bots_backend/internal/repo/api_integrations"
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
	// TODO: move postgres init

	link := "postgres://%s:%s@%s:%d/%s"

	postgresCFG, err := pgxpool.ParseConfig(fmt.Sprintf(link,
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name))
	if err != nil {
		return nil, err
	}

	postgresCFG.MaxConns = 100
	postgresCFG.MinConns = 0
	postgresCFG.MaxConnLifetime = 1 * time.Hour
	postgresCFG.MaxConnIdleTime = 10 * time.Minute

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), postgresCFG)
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
	go func() {
		// TODO: CONSUMER
		err := TestRead()
		if err != nil {
			fmt.Printf("ERRROR!!!!!!!, %s\n", err.Error())
		}
	}()
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
