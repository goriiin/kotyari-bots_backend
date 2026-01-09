package posts_command_producer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/ierrors"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandHandler) CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error) {
	bot, err := p.fetcher.GetBot(ctx, req.BotId.String())
	if err != nil {
		p.log.Error(err, true, "CreatePost: get bot")
		return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: ierrors.GRPCToDomainError(err).Error()}, nil
	}

	idsString := make([]string, 0, len(req.ProfileIds))
	for _, id := range req.ProfileIds {
		idsString = append(idsString, id.String())
	}

	profilesBatch, err := p.fetcher.GetProfiles(ctx, idsString)
	if err != nil {
		p.log.Error(err, true, "CreatePost: get profiles")
		return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: ierrors.GRPCToDomainError(err).Error()}, nil
	}

	postProfiles := make([]posts.CreatePostProfiles, 0, len(idsString))
	for _, profile := range profilesBatch.Profiles {
		profileID, _ := uuid.Parse(profile.Id)
		postProfiles = append(postProfiles, posts.CreatePostProfiles{
			ProfileID:     profileID,
			ProfilePrompt: profile.Prompt,
			ProfileName:   profile.Name,
		})
	}

	groupID := uuid.New()
	botID, _ := uuid.Parse(bot.Id)
	createPostRequest := posts.KafkaCreatePostRequest{
		PostID:             uuid.New(),
		GroupID:            groupID,
		BotID:              botID,
		BotName:            bot.BotName,
		BotPrompt:          bot.BotPrompt,
		UserPrompt:         req.TaskText,
		Profiles:           postProfiles,
		Platform:           model.PlatformType(req.Platform),
		PostType:           model.PostType(req.PostType.Value),
		ModerationRequired: bot.ModerationRequired,
	}

	rawReq, err := jsoniter.Marshal(createPostRequest)
	if err != nil {
		p.log.Error(err, true, "CreatePost: marshal")
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdCreate, createPostRequest.GroupID.String(), rawReq), 30*time.Second)
	if err != nil {
		p.log.Error(err, true, "CreatePost: request")
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		p.log.Error(err, true, "CreatePost: unmarshal response")
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	if resp.Error != "" {
		p.log.Warn("CreatePost: response error", errors.New(resp.Error))
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   fmt.Sprintf("Failed to create post, %s", resp.Error),
		}, nil
	}

	return &gen.PostCreateResponse{GroupID: groupID}, nil
}
