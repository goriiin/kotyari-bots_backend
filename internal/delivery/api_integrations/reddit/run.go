package reddit

import (
	"context"
	"time"
)

func (r *RedditAPIDelivery) Run() error {
	runRequest := func() error {
		news, err := r.performRequests()
		if err != nil {
			// TODO: log, err
			return err
		}
		for _, redditNews := range news {
			for _, post := range redditNews.Data.Posts {
				err := r.producer.Publish(context.Background(), []byte(post.PostData.Title))
				if err != nil {
					return err
				}

				//fmt.Println(post.PostData.Title)
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
