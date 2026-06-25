package posts_command_producer

import (
	"context"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandHandler) CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return nil, err
	}

	bot, err := p.fetcher.GetBot(ctx, req.BotId.String())
	if err != nil {
		p.log.Error(err, true, "CreatePost: get bot")
		return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: constants.InternalMsg}, nil
	}

	idsString := make([]string, 0, len(req.ProfileIds))
	for _, id := range req.ProfileIds {
		idsString = append(idsString, id.String())
	}

	profilesBatch, err := p.fetcher.GetProfiles(ctx, idsString)
	if err != nil {
		p.log.Error(err, true, "CreatePost: get profiles")
		return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: constants.InternalMsg}, nil
	}

	postProfiles := make([]posts.CreatePostProfiles, 0, len(idsString))
	for _, profile := range profilesBatch.Profiles {
		profileID, parseErr := uuid.Parse(profile.Id)
		if parseErr != nil {
			p.log.Error(parseErr, true, "CreatePost: parse profile id")
			return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: constants.InternalMsg}, nil
		}
		postProfiles = append(postProfiles, posts.CreatePostProfiles{
			ProfileID:     profileID,
			ProfilePrompt: profile.Prompt,
			ProfileName:   profile.Name,
		})
	}

	groupID := uuid.New()
	botID, parseErr := uuid.Parse(bot.Id)
	if parseErr != nil {
		p.log.Error(parseErr, true, "CreatePost: parse bot id")
		return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: constants.InternalMsg}, nil
	}
	createPostRequest := posts.KafkaCreatePostRequest{
		PostID:             uuid.New(),
		UserID:             userID,
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
			Message:   constants.InternalMsg,
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdCreate, createPostRequest.GroupID.String(), rawReq), 30*time.Second)
	if err != nil {
		p.log.Error(err, true, "CreatePost: request")
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		p.log.Error(err, true, "CreatePost: unmarshal response")
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	if resp.Error != "" {
		p.log.Warn("CreatePost: response error", errors.New(resp.Error))
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	return &gen.PostCreateResponse{GroupID: groupID}, nil
}
