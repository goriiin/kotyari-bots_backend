package posts_command_consumer

import (
	"context"
	"fmt"
	"log"
	"sync"

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
		go func(profile posts.CreatePostProfiles) {
			defer wg.Done()
			post := p.processProfile(ctx, req, profile, postsMap)
			if post != nil {
				postsChan <- *post
			}
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
		return errors.Wrap(err, "failed to create posts")
	}

	return nil
}

// processProfile processes a single profile and returns the best post
func (p *PostsCommandConsumer) processProfile(ctx context.Context, req posts.KafkaCreatePostRequest, profile posts.CreatePostProfiles, postsMap map[uuid.UUID]model.Post) *model.Post {
	profilesPosts := p.generatePostsForProfile(ctx, req, profile, postsMap)
	if len(profilesPosts) == 0 {
		return nil
	}

	bestPostCandidate, err := p.judge.SelectBest(ctx, req.UserPrompt, profile.ProfilePrompt, req.BotPrompt,
		posts.PostsToCandidates(profilesPosts))
	if err != nil {
		fmt.Println("error getting best post ", err)
		return nil
	}

	bestPost := p.createPostFromCandidate(req, profile, bestPostCandidate)
	p.publishToOtvet(ctx, req, bestPostCandidate, bestPost)

	return bestPost
}

// generatePostsForProfile generates multiple post candidates for a profile
func (p *PostsCommandConsumer) generatePostsForProfile(ctx context.Context, req posts.KafkaCreatePostRequest, profile posts.CreatePostProfiles, postsMap map[uuid.UUID]model.Post) []model.Post {
	rewritten, err := p.rewriter.Rewrite(ctx, req.UserPrompt, profile.ProfilePrompt, req.BotPrompt)
	if err != nil {
		fmt.Println("error rewriting prompts", err)
		return nil
	}

	var (
		mutex     sync.Mutex
		profileWg sync.WaitGroup
	)

	profilesPosts := make([]model.Post, 0, len(rewritten))

	for _, rw := range rewritten {
		profileWg.Add(1)
		go func(rewrittenPrompt string) {
			defer profileWg.Done()

			generatedPostContent, err := p.getter.GetPost(ctx, rewrittenPrompt, profile.ProfilePrompt, req.BotPrompt)
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
		}(rw)
	}

	profileWg.Wait()
	return profilesPosts
}


// createPostFromCandidate creates a Post model from a candidate
func (p *PostsCommandConsumer) createPostFromCandidate(req posts.KafkaCreatePostRequest, profile posts.CreatePostProfiles, candidate model.Candidate) *model.Post {
	return &model.Post{
		ID:          uuid.New(),
		OtvetiID:    0,
		BotID:       req.BotID,
		BotName:     req.BotName,
		ProfileID:   profile.ProfileID,
		ProfileName: profile.ProfileName,
		GroupID:     req.GroupID,
		Platform:    req.Platform,
		Type:        req.PostType,
		UserPrompt:  req.UserPrompt,
		Title:       candidate.Title,
		Text:        candidate.Text,
	}
}


// publishToOtvet publishes post to otvet.mail.ru if platform is otveti
func (p *PostsCommandConsumer) publishToOtvet(ctx context.Context, req posts.KafkaCreatePostRequest, candidate model.Candidate, post *model.Post) {
	if req.Platform != model.OtvetiPlatform || p.otvetClient == nil {
		return
	}

	topicType := getTopicTypeFromPostType(req.PostType)
	spaces := p.getSpacesForPost(ctx, candidate)

	log.Printf("INFO: topicType: %+v\t spaces: %+v\n", topicType, spaces)

	otvetResp, err := p.otvetClient.CreatePostSimple(ctx, candidate.Title, candidate.Text, topicType, spaces)
	if err != nil {
		fmt.Printf("error publishing post to otvet: %v\n", err)
		return
	}

	log.Printf("INFO: published post to otvet: %v\n", otvetResp)

	if otvetResp != nil && otvetResp.Result != nil {
		post.OtvetiID = uint64(otvetResp.Result.ID)
	}
}


// getSpacesForPost predicts spaces for a post or returns default spaces
func (p *PostsCommandConsumer) getSpacesForPost(ctx context.Context, candidate model.Candidate) []otvet.Space {
	combinedText := candidate.Title + " " + candidate.Text
	spaces := getDefaultSpaces()

	predictResp, err := p.otvetClient.PredictTagsSpaces(ctx, combinedText)
	if err != nil {
		fmt.Printf("error predicting spaces: %v, using default spaces\n", err)
		return spaces
	}

	if predictResp == nil || len(*predictResp) == 0 {
		return spaces
	}

	// Convert predicted spaces to Space format
	predictedSpaces := make([]otvet.Space, 0, len((*predictResp)[0].Spaces))
	for _, spaceID := range (*predictResp)[0].Spaces {
		predictedSpaces = append(predictedSpaces, otvet.Space{
			ID:      spaceID,
			IsPrime: true,
		})
	}

	if len(predictedSpaces) > 0 {
		return predictedSpaces
	}

	return spaces
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
