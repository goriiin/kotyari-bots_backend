package aggregator

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/utils"
	"github.com/segmentio/kafka-go"
)

// TODO: remove sources
const redditSource = "reddit"

func (s *AggregatorService) AddTopics(ctx context.Context, messages []kafka.Message) error {
	topics := make([]model.Topic, 0, len(messages))

	for _, message := range messages {
		messageText := string(message.Value)
		messageHash := utils.HashString(messageText)

		topic := model.Topic{
			Source: redditSource,
			Text:   messageText,
			Hash:   messageHash,
		}

		topics = append(topics, topic)
	}

	if err := s.topicsCreator.AddTopics(ctx, topics); err != nil {
		return err
	}

	return nil
}
