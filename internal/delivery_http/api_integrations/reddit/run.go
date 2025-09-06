package reddit

import (
	"context"
	"time"
)

func (r *RedditAPIDelivery) Run() error {
	runRequest := func() error {
		postChan, err := r.performRequests()
		if err != nil {
			return err
		}

		for post := range postChan {
			err := r.producer.Publish(context.Background(), []byte(post.Title))
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := runRequest(); err != nil {
		r.log.Error().Stack().Err(err).Msg("error happened while performing request")
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := runRequest(); err != nil {
			r.log.Error().Stack().Err(err).Msg("error happened while performing request")
		}
	}

	return nil
}
