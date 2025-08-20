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

func (r *RedditAPIDelivery) performRequests() ([]RedditAPIResponse, error) {
	// TODO: log

	ctx, cancel := context.WithTimeout(context.Background(), defaultErrGroupWaitTime)
	defer cancel()

	integrations, err := r.integration.GetIntegrations(ctx, redditAPIString)
	if err != nil {
		// TODO: log, err
		return nil, err
	}

	redditAPIResponses := make([]RedditAPIResponse, 0, len(integrations))
	g, _ := errgroup.WithContext(ctx)

	var mu sync.Mutex
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

			mu.Lock()
			redditAPIResponses = append(redditAPIResponses, redditAPIResponse)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return redditAPIResponses, err
	}

	return redditAPIResponses, nil
}
