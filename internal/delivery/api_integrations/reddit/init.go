package reddit

import (
	"context"
	"net/http"
	"time"

	"github.com/kotyari-bots_backend/internal/model"
)

type IntegrationsRepo interface {
	GetIntegrations(ctx context.Context, integrationName string) ([]model.APIIntegration, error)
}

type MessageProducer interface {
	Publish(ctx context.Context, message []byte) error
}

type RedditAPIDelivery struct {
	client      *http.Client
	integration IntegrationsRepo
	interval    time.Duration
	producer    MessageProducer
}

func NewRedditApiIntegration(client *http.Client, repo IntegrationsRepo, interval time.Duration, producer MessageProducer) *RedditAPIDelivery {
	return &RedditAPIDelivery{
		client:      client,
		integration: repo,
		interval:    interval,
		producer:    producer,
	}
}
