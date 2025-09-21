package reddit

import (
	"fmt"
	"time"
)

func (r *RedditAPIDelivery) Run() error {
	runRequest := func() error {
		postChan, err := r.performRequests()
		if err != nil {
			// TODO: log, err
			return err
		}

		for post := range postChan {
			// err := r.producer.Publish(context.Background(), []byte(post.Title))
			fmt.Println(post.Comments)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := runRequest(); err != nil {
		return err
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := runRequest(); err != nil {
			return err
		}
	}

	return nil
}
