package posts_command_consumer

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandConsumer) CreatePost(ctx context.Context, payload []byte) error {
	var req posts.KafkaCreatePostRequest
	err := jsoniter.Unmarshal(payload, &req)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	//postsChan := make(chan model.Post, len(req.Profiles))
	//var wg sync.WaitGroup
	//for _, profile := range req.Profiles {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//
	//		generatedPostContent, err := p.getter.GetPost(ctx, req.UserPrompt, profile.ProfilePrompt, req.BotPrompt)
	//		if err != nil {
	//			fmt.Println("error getting post", err)
	//			return
	//		}
	//
	//		post := model.Post{
	//			ID:         req.PostID,
	//			OtvetiID:   0, // Пока так
	//			BotID:      req.BotID,
	//			ProfileID:  profile.ProfileID,
	//			GroupID:    req.GroupID,
	//			Platform:   req.Platform,
	//			Type:       req.PostType,
	//			UserPrompt: req.UserPrompt,
	//			Title:      generatedPostContent.PostTitle,
	//			Text:       generatedPostContent.PostText,
	//		}
	//
	//		postsChan <- post
	//	}()
	//}
	//
	//go func() {
	//	wg.Wait()
	//	close(postsChan)
	//}()

	postTest := []model.Post{
		{
			ID:         uuid.New(),
			OtvetiID:   0, // Пока так
			BotID:      req.BotID,
			ProfileID:  req.Profiles[0].ProfileID,
			GroupID:    req.GroupID,
			Platform:   req.Platform,
			Type:       req.PostType,
			UserPrompt: req.UserPrompt,
			Title:      "TITLE ONE",
			Text:       "TEXT ONE",
		},
		{
			ID:         uuid.New(),
			OtvetiID:   0, // Пока так
			BotID:      req.BotID,
			ProfileID:  req.Profiles[1].ProfileID,
			GroupID:    req.GroupID,
			Platform:   req.Platform,
			Type:       req.PostType,
			UserPrompt: req.UserPrompt,
			Title:      "TITLE ZWEI",
			Text:       "TEXT ZWEI",
		},
	}

	//finalPosts := make([]model.Post, 0, len(req.Profiles))
	//for post := range postsChan {
	//	finalPosts = append(finalPosts, post)
	//}

	err = p.repo.CreatePostsBatch(ctx, postTest)
	if err != nil {
		return errors.Wrap(err, "failed to create posts")
	}

	return nil
}
