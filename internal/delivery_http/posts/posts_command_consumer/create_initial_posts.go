package posts_command_consumer

import (
	"context"
	"maps"
	"slices"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandConsumer) CreateInitialPosts(ctx context.Context, payload []byte) (map[uuid.UUID]model.Post, posts.KafkaCreatePostRequest, error) {
	var req posts.KafkaCreatePostRequest
	err := jsoniter.Unmarshal(payload, &req)
	if err != nil {
		return nil, posts.KafkaCreatePostRequest{}, errors.Wrapf(constants.ErrInternal, "failed to unmarshal: %s", err.Error())
	}

	initialPosts := make(map[uuid.UUID]model.Post, len(req.Profiles))

	for _, profile := range req.Profiles {
		post := model.Post{
			ID:          uuid.New(),
			UserID:      req.UserID, // Map UserID from Kafka Request
			OtvetiID:    0,
			BotID:       req.BotID,
			BotName:     req.BotName,
			ProfileID:   profile.ProfileID,
			ProfileName: profile.ProfileName,
			GroupID:     req.GroupID,
			Platform:    req.Platform,
			Type:        req.PostType,
			UserPrompt:  req.UserPrompt,
			Title:       "",
			Text:        "",
		}

		initialPosts[profile.ProfileID] = post
	}

	err = p.repo.CreatePostsBatch(ctx, slices.Collect(maps.Values(initialPosts)))
	if err != nil {
		return nil, posts.KafkaCreatePostRequest{}, errors.Wrapf(err, "failed to create initial posts")
	}

	return initialPosts, req, nil
}
