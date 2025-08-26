package aggregator

import (
	"context"
	"sync"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/utils"
	"github.com/segmentio/kafka-go"
)

// TODO: remove sources
const redditSource = "reddit"

func (s *AggregatorService) AddTopics(ctx context.Context, messages []kafka.Message) error {
	topicsChan := make(chan model.Topic, len(messages))
	var wg sync.WaitGroup

	for _, message := range messages {
		wg.Add(1)
		go func(message kafka.Message) {
			defer wg.Done()
			messageHash := utils.HashString(message.Value)
			topic := model.Topic{
				Source: redditSource,
				Text:   string(message.Value),
				Hash:   messageHash,
			}
			topicsChan <- topic
		}(message)
	}

	wg.Wait()
	close(topicsChan)

	topics := make([]model.Topic, 0, len(messages))
	for topic := range topicsChan {
		topics = append(topics, topic)
	}

	if err := s.topicsCreator.AddTopics(ctx, topics); err != nil {
		return err
	}

	return nil
}
