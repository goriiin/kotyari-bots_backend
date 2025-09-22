package reddit

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"golang.org/x/sync/errgroup"
)

// TODO: move in future
const (
	redditAPIString         = "reddit"
	defaultErrGroupWaitTime = 20 * time.Second
)

func (r *RedditAPIDelivery) performRequests() (chan PostData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultErrGroupWaitTime)
	defer cancel()

	integrations, err := r.integration.GetIntegrations(ctx, redditAPIString)
	if err != nil {
		return nil, err
	}

	redditAPIResponses := make(chan RedditAPIResponse)
	g, _ := errgroup.WithContext(ctx)

	for _, integration := range integrations {
		g.Go(func() error {
			req, err := http.NewRequest(http.MethodGet, integration.Url, http.NoBody)
			if err != nil {
				return errors.Wrap(err, "failed to create request")
			}

			resp, err := r.client.Do(req)
			if err != nil {
				return errors.Wrap(err, "failed to perform request")
			}

			if resp.StatusCode == http.StatusForbidden {
				return errors.New("request was blocked")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrapf(err, "bad response body: %s", string(body))
			}

			err = resp.Body.Close()
			if err != nil {
				return errors.Wrap(err, "failed to close body")
			}

			var redditAPIResponse RedditAPIResponse
			if err := json.Unmarshal(body, &redditAPIResponse); err != nil {
				return errors.Wrapf(err, "failed to unmarhsal: %s", integration.Url)
			}
			redditAPIResponses <- redditAPIResponse

			return nil
		})
	}

	go func() {
		defer close(redditAPIResponses)
		if err := g.Wait(); err != nil {
			// TODO: add error behaviour
			r.log.Error(err, false, "failed wait")
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
