package posts_command_consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

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
		go func(prof posts.CreatePostProfiles) {
			defer wg.Done()

			profileCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
			defer cancel()

			var mutex sync.Mutex
			profileWg := sync.WaitGroup{}
			profilesPosts := make([]model.Post, 0, 5)

			rewritten, err := p.rewriter.Rewrite(profileCtx, req.UserPrompt, prof.ProfilePrompt, req.BotPrompt)
			if err != nil {
				fmt.Println("error rewriting prompts:", err)
				return
			}

			for _, rw := range rewritten {
				profileWg.Add(1)
				go func(rewrittenText string) {
					defer profileWg.Done()

					generatedPostContent, err := p.getter.GetPost(profileCtx, rewrittenText, prof.ProfilePrompt, req.BotPrompt)
					if err != nil {
						fmt.Println("error getting post:", err)
						return
					}

					post := postsMap[prof.ProfileID]
					post.Title = generatedPostContent.PostTitle
					post.Text = generatedPostContent.PostText

					mutex.Lock()
					profilesPosts = append(profilesPosts, post)
					mutex.Unlock()
				}(rw)
			}

			profileWg.Wait()

			bestPostCandidate, err := p.judge.SelectBest(profileCtx, req.UserPrompt, prof.ProfilePrompt, req.BotPrompt, posts.PostsToCandidates(profilesPosts))
			if err != nil {
				fmt.Println("error getting best post:", err)
				return
			}

			bestPost := postsMap[prof.ProfileID]
			bestPost.Text = bestPostCandidate.Text
			bestPost.Title = bestPostCandidate.Title

			postsChan <- bestPost
		}(profile)
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
