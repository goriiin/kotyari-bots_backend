package aggregator

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type TopicsCreator interface {
	AddTopics(ctx context.Context, topics []model.Topic) error
}

type AggregatorService struct {
	log           *logger.Logger
	topicsCreator TopicsCreator
}

func NewAggregatorService(log *logger.Logger, creator TopicsCreator) *AggregatorService {
	return &AggregatorService{topicsCreator: creator, log: log}
}
