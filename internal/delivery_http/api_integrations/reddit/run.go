package reddit

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

func (r *RedditAPIDelivery) Run() error {
	runRequest := func() error {
		postChan, err := r.performRequests()
		if err != nil {
			return err
		}

		for post := range postChan {
			err := r.producer.Publish(context.Background(), kafka.Message{Value: []byte(post.Title)})
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := runRequest(); err != nil {
		r.log.Error(err, true, "error happened while performing request")
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := runRequest(); err != nil {
			r.log.Error(err, true, "error happened while performing request")
		}
	}

	return nil
}
