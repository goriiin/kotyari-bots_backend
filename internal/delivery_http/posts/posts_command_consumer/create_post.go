package posts_command_consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (p *PostsCommandConsumer) CreatePost(ctx context.Context, postsMap map[uuid.UUID]model.Post, req posts.KafkaCreatePostRequest) error {
	postsChan := make(chan model.Post, len(req.Profiles))
	var wg sync.WaitGroup
	for _, profile := range req.Profiles {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var (
				mutex     sync.Mutex
				profileWg sync.WaitGroup
			)

			profilesPosts := make([]model.Post, 0, 5)

			rewritten, err := p.rewriter.Rewrite(ctx, req.UserPrompt, profile.ProfilePrompt, req.BotPrompt)
			if err != nil {
				fmt.Println("error rewriting prompts", err)
				return
			}

			for _, rw := range rewritten {
				profileWg.Add(1)
				go func() {
					defer profileWg.Done()

					generatedPostContent, err := p.getter.GetPost(ctx, rw, profile.ProfilePrompt, req.BotPrompt)
					if err != nil {
						fmt.Println("error getting post", err)
						return
					}

					post := postsMap[profile.ProfileID]
					post.Title = generatedPostContent.PostTitle
					post.Text = generatedPostContent.PostText

					mutex.Lock()
					profilesPosts = append(profilesPosts, post)
					mutex.Unlock()
				}()
			}

			profileWg.Wait()

			fmt.Println("PROFILE_POSTS: ", profilesPosts)

			bestPostCandidate, err := p.judge.SelectBest(ctx, req.UserPrompt, profile.ProfilePrompt, req.BotPrompt,
				posts.PostsToCandidates(profilesPosts))

			fmt.Println("TITLE: ", bestPostCandidate.Title, "TEXT: ", bestPostCandidate.Text)

			if err != nil {
				fmt.Println("error getting best post ", err)
				return
			}

			bestPost := postsMap[profile.ProfileID]

			bestPost.Text = bestPostCandidate.Text
			bestPost.Title = bestPostCandidate.Title

			postsChan <- bestPost
		}()
	}

	go func() {
		wg.Wait()
		close(postsChan)
	}()

	finalPosts := make([]model.Post, 0, len(req.Profiles))
	for post := range postsChan {
		finalPosts = append(finalPosts, post)
	}

	err := p.repo.UpdatePostsBatch(ctx, finalPosts)
	if err != nil {
		return errors.Wrap(err, "failed to update posts")
	}

	return nil
}
