package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/sync/errgroup"
)

// TODO: move in future
const (
	redditAPIString         = "reddit"
	defaultErrGroupWaitTime = 120 * time.Second
)

func (r *RedditAPIDelivery) performRequests() (chan FullPost, error) {
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

			if resp.StatusCode == http.StatusTooManyRequests {
				fmt.Println("too many requests O GORE DRANOYE")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// TODO: log, err

				return fmt.Errorf("bad resposeBody: %w\n%s", err, string(body))
			}

			//fmt.Printf("body: %s\n", string(body))

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
			//fmt.Println("resp:", redditAPIResponse)

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

	posts := make(chan FullPost)

	go func() {
		defer close(posts)

		postsErrG, _ := errgroup.WithContext(ctx)

		for resp := range redditAPIResponses {
			for _, post := range resp.Data.Posts {
				postsErrG.Go(func() error {
					comments, err := r.fetchComments(post.PostData)
					if err != nil {
						// TODO: log

						fmt.Printf("comments uvi, %s\n", err.Error())
					}

					fullPost := FullPost{
						Post:     post.PostData,
						Comments: comments,
					}

					posts <- fullPost

					return nil
				})
			}
		}
	}()

	//	var wg sync.WaitGroup
	//	go func() {
	//		for redditNews := range redditAPIResponses {
	//			wg.Add(1)
	//			for _, post := range redditNews.Data.Posts {
	//				posts <- post.PostData
	//			}
	//			wg.Done()
	//		}
	//		wg.Wait()
	//		close(posts)
	//	}()
	//
	//	return posts, nil

	return posts, err
}

func (r *RedditAPIDelivery) fetchComments(post PostData) ([]CommentData, error) {
	commentsURL := fmt.Sprintf("https://www.reddit.com/r/%s/comments/%s.json", url.PathEscape(post.Subreddit), url.PathEscape(post.ID))

	req, err := http.NewRequest(http.MethodGet, commentsURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create comments request: %w", err)
	}
	time.Sleep(20 * time.Second)
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		fmt.Println("429 for dranie commenti O GORE BALIN(((((")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read comments body: %w", err)
	}

	var commentResponse CommentAPIResponse
	if err := json.Unmarshal(body, &commentResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal comments: %w", err)
	}

	if len(commentResponse) < 2 {
		return []CommentData{}, nil
	}

	commentListings := commentResponse[1].Data.Children
	comments := make([]CommentData, 0, len(commentListings))
	for _, comment := range commentListings {
		if comment.Data.Author == "[deleted]" || comment.Data.Body == "" {
			continue
		}
		comments = append(comments, CommentData{
			Author: comment.Data.Author,
			Body:   comment.Data.Body,
			Score:  comment.Data.Score,
		})
	}

	return comments, nil
}
