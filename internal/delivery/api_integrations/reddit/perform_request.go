package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// TODO: move in future
const (
	redditAPIString         = "reddit"
	defaultErrGroupWaitTime = 20 * time.Second
)

func (r *RedditAPIDelivery) performRequests() (chan PostData, error) {
	// TODO: log

	ctx, cancel := context.WithTimeout(context.Background(), defaultErrGroupWaitTime)
	defer cancel()

	integrations, err := r.integration.GetIntegrations(ctx, redditAPIString)
	if err != nil {
		// TODO: log, err
		return nil, err
	}

	redditAPIResponses := make(chan RedditAPIResponse)
	g, _ := errgroup.WithContext(ctx)

	for _, integration := range integrations {
		g.Go(func() error {
			req, err := http.NewRequest(http.MethodGet, integration.Url, http.NoBody)
			if err != nil {
				// TODO: log, err

				return fmt.Errorf("failed to create request: %w", err)
			}

			resp, err := r.client.Do(req)
			if err != nil {
				// TODO: log, err

				return fmt.Errorf("failed to do request: %w", err)
			}

			// TODO: add resp.StatusCode check

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// TODO: log, err

				return fmt.Errorf("bad resposeBody: %w\n%s", err, string(body))
			}

			fmt.Printf("body: %s\n", string(body))

			err = resp.Body.Close()
			if err != nil {
				// TODO: log, err

				return fmt.Errorf("failed to close body: %w", err)
			}

			var redditAPIResponse RedditAPIResponse
			if err := json.Unmarshal(body, &redditAPIResponse); err != nil {
				// TODO: log, err

				return fmt.Errorf("failed to unmarshal: %s", integration.Url)
			}
			redditAPIResponses <- redditAPIResponse
			fmt.Println("resp:", redditAPIResponse)

			return nil
		})
	}

	go func() {
		defer close(redditAPIResponses)
		if err := g.Wait(); err != nil {
			// TODO: add error behaviour
			fmt.Println("uvi")
		}
	}()

	posts := make(chan PostData)

	var wg sync.WaitGroup
	go func() {
		for redditNews := range redditAPIResponses {
			wg.Add(1)
			for _, post := range redditNews.Data.Posts {
				posts <- post.PostData
			}
			wg.Done()
		}
		wg.Wait()
		close(posts)
	}()

	return posts, nil
}
