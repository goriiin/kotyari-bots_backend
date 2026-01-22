package posts_command_consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/otvet"
	"github.com/goriiin/kotyari-bots_backend/pkg/posting_queue"
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
		p.log.Error(err, true, "CreatePost: update posts batch")
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
		p.log.Error(err, true, "processProfile: select best")
		return nil
	}

	bestPost := postsMap[profile.ProfileID]
	bestPost.Title = bestPostCandidate.Title
	bestPost.Text = bestPostCandidate.Text

	p.publishToOtvet(ctx, req, bestPostCandidate, &bestPost)

	return &bestPost
}

// generatePostsForProfile generates multiple post candidates for a profile
func (p *PostsCommandConsumer) generatePostsForProfile(ctx context.Context, req posts.KafkaCreatePostRequest, profile posts.CreatePostProfiles, postsMap map[uuid.UUID]model.Post) []model.Post {
	rewritten, err := p.rewriter.Rewrite(ctx, req.UserPrompt, profile.ProfilePrompt, req.BotPrompt)
	if err != nil {
		p.log.Error(err, true, "generatePostsForProfile: rewrite")
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
				p.log.Error(err, true, "generatePostsForProfile: get post")
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

// publishToOtvet publishes post to otvet.mail.ru if platform is otveti
func (p *PostsCommandConsumer) publishToOtvet(ctx context.Context, req posts.KafkaCreatePostRequest, candidate model.Candidate, post *model.Post) {
	if req.Platform != model.OtvetiPlatform || p.otvetClient == nil {
		return
	}

	topicType := getTopicTypeFromPostType(req.PostType)
	spaces := p.getSpacesForPost(ctx, candidate)

	if p.queue != nil {
		postRequest := posting_queue.PostRequest{
			Platform:           req.Platform,
			PostType:           req.PostType,
			TopicType:          topicType,
			Spaces:             spaces,
			ModerationRequired: req.ModerationRequired,
		}
		p.queue.Enqueue(post, candidate, postRequest)
		p.log.Info(fmt.Sprintf("publishToOtvet: post %s enqueued", post.ID.String()))
		return
	}

	if req.ModerationRequired {
		p.log.Info(fmt.Sprintf("publishToOtvet: bot requires moderation, skipping direct publish for post %s", post.ID.String()))
		return
	}

	p.log.Info(fmt.Sprintf("publishToOtvet: attempting direct publish for post %s", post.ID.String()))
	otvetResp, err := p.otvetClient.CreatePostSimple(ctx, candidate.Title, candidate.Text, topicType, spaces)
	if err != nil {
		p.log.Error(err, true, "publishToOtvet: create post simple failed")
		return
	}

	p.log.Info(fmt.Sprintf("publishToOtvet: success for post %s, response: %+v", post.ID.String(), otvetResp))

	if otvetResp != nil && otvetResp.Result != nil {
		post.OtvetiID = uint64(otvetResp.Result.ID)
		post.IsPublished = true
		post.URL = fmt.Sprintf("https://otvet.mail.ru/question/%d", otvetResp.Result.ID)
	}
}

// getSpacesForPost predicts spaces for a post or returns default spaces
func (p *PostsCommandConsumer) getSpacesForPost(ctx context.Context, candidate model.Candidate) []otvet.Space {
	combinedText := candidate.Title + " " + candidate.Text
	spaces := getDefaultSpaces()

	predictResp, err := p.otvetClient.PredictTagsSpaces(ctx, combinedText)
	if err != nil {
		p.log.Error(err, true, "getSpacesForPost: predict tags spaces")
		return spaces
	}

	if predictResp == nil || len(*predictResp) == 0 {
		return spaces
	}

	predictedSpaces := make([]otvet.Space, 0, len((*predictResp)[0].Spaces))
	for i, spaceID := range (*predictResp)[0].Spaces {
		predictedSpaces = append(predictedSpaces, otvet.Space{
			ID:      spaceID,
			IsPrime: i == 0,
		})
	}

	if len(predictedSpaces) > 0 {
		return predictedSpaces
	}

	return spaces
}

func getTopicTypeFromPostType(postType model.PostType) int {
	switch postType {
	case model.OpinionPostType:
		return 2
	case model.KnowledgePostType:
		return 2
	case model.HistoryPostType:
		return 2
	default:
		return 2
	}
}

func getDefaultSpaces() []otvet.Space {
	return []otvet.Space{
		{
			ID:      501,
			IsPrime: true,
		},
	}
}
