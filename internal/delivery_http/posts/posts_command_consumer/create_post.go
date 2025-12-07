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
	"github.com/goriiin/kotyari-bots_backend/pkg/otvet"
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

			// Publish to otvet.mail.ru if platform is otveti
			if req.Platform == model.OtvetiPlatform && p.otvetClient != nil {
				topicType := getTopicTypeFromPostType(req.PostType)

				// Predict spaces from title + text
				combinedText := bestPostCandidate.Title + " " + bestPostCandidate.Text
				spaces := getDefaultSpaces() // fallback to default

				predictResp, err := p.otvetClient.PredictTagsSpaces(ctx, combinedText)
				if err != nil {
					fmt.Printf("error predicting spaces: %v, using default spaces\n", err)
				} else if predictResp != nil && len(*predictResp) > 0 {
					// Convert predicted spaces to Space format
					predictedSpaces := make([]otvet.Space, 0, len((*predictResp)[0].Spaces))
					for _, spaceID := range (*predictResp)[0].Spaces {
						predictedSpaces = append(predictedSpaces, otvet.Space{
							ID:      spaceID,
							IsPrime: true, // Default value, can be adjusted if needed
						})
					}
					if len(predictedSpaces) > 0 {
						spaces = predictedSpaces
					}
				}

				otvetResp, err := p.otvetClient.CreatePostSimple(ctx, bestPostCandidate.Title, bestPostCandidate.Text, topicType, spaces)
				if err != nil {
					fmt.Printf("error publishing post to otvet: %v\n", err)
					// Continue anyway, post will be saved without OtvetiID
				} else if otvetResp != nil && otvetResp.Result != nil {
					bestPost.OtvetiID = uint64(otvetResp.Result.ID)
				}
			}

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

// getTopicTypeFromPostType converts PostType to otvet topic_type
// topic_type: 2 = question (opinion), other values may be used for other types
func getTopicTypeFromPostType(postType model.PostType) int {
	switch postType {
	case model.OpinionPostType:
		return 2 // question
	case model.KnowledgePostType:
		return 2 // question (can be adjusted if needed)
	case model.HistoryPostType:
		return 2 // question (can be adjusted if needed)
	default:
		return 2 // default to question
	}
}

// getDefaultSpaces returns default spaces for otvet posts
// TODO: move to config or get from request
func getDefaultSpaces() []otvet.Space {
	// Default space - can be configured later
	return []otvet.Space{
		{
			ID:      501, // Example space ID from the response
			IsPrime: true,
		},
	}
}
