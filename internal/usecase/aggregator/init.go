package aggregator

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type TopicsCreator interface {
	AddTopics(ctx context.Context, topics []model.Topic) error
}

type AggregatorService struct {
	// TODO: logs

	topicsCreator TopicsCreator
}

func NewAggregatorService(creator TopicsCreator) *AggregatorService {
	return &AggregatorService{topicsCreator: creator}
}
