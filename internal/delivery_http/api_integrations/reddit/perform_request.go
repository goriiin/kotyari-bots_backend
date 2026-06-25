package reddit

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
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
	// NOTE: ownership of cancel is handed to the draining goroutine below — we
	// must NOT defer-cancel here, or the context would be cancelled the moment
	// this function returns, killing the goroutines that still feed `posts`.
	ctx, cancel := context.WithTimeout(context.Background(), defaultErrGroupWaitTime)

	integrations, err := r.integration.GetIntegrations(ctx, redditAPIString)
	if err != nil {
		cancel()
		return nil, err
	}

	redditAPIResponses := make(chan RedditAPIResponse)
	g, _ := errgroup.WithContext(ctx)

	for _, integration := range integrations {
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, integration.Url, http.NoBody)
			if err != nil {
				return errors.Wrap(err, "failed to create request")
			}

			resp, err := r.client.Do(req)
			if err != nil {
				return errors.Wrap(err, "failed to perform request")
			}
			// Always close the body, including on the Forbidden early-return
			// below (previously leaked the connection on every blocked request).
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode == http.StatusForbidden {
				return errors.New("request was blocked")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrapf(err, "bad response body: %s", string(body))
			}

			var redditAPIResponse RedditAPIResponse
			if err := json.Unmarshal(body, &redditAPIResponse); err != nil {
				return errors.Wrapf(err, "failed to unmarhsal: %s", integration.Url)
			}

			// Respect cancellation so a stalled consumer can't wedge this
			// goroutine forever on an unbuffered send.
			select {
			case redditAPIResponses <- redditAPIResponse:
			case <-ctx.Done():
				return ctx.Err()
			}

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

	go func() {
		defer cancel()
		defer close(posts)
		for redditNews := range redditAPIResponses {
			for _, post := range redditNews.Data.Posts {
				select {
				case posts <- post.PostData:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return posts, nil
}
