package posts_command_producer

import (
	"context"
	"net/http"

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
		p.log.Error(err, true, "failed to get bot info")
		return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: ierrors.GRPCToDomainError(err).Error()}, nil
	}

	idsString := make([]string, 0, len(req.ProfileIds))
	for _, id := range req.ProfileIds {
		idsString = append(idsString, id.String())
	}

	profilesBatch, err := p.fetcher.GetProfiles(ctx, idsString)
	if err != nil {
		p.log.Error(err, true, "failed to get profiles info")
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
	createPostRequest := posts.KafkaCreatePostRequest{
		GroupID:    groupID,
		BotID:      req.BotId,
		BotName:    bot.BotName,
		BotPrompt:  bot.BotPrompt,
		UserPrompt: req.TaskText,
		Profiles:   postProfiles,
		Platform:   model.PlatformType(req.Platform),
		PostType:   model.PostType(req.PostType.Value),
	}

	rawReq, err := jsoniter.Marshal(createPostRequest)
	if err != nil {
		p.log.Error(err, true, "failed to marshal create post request")
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	err = p.producer.Publish(ctx, posts.PayloadToEnvelope(posts.CmdCreate, createPostRequest.GroupID.String(), rawReq))
	if err != nil {
		p.log.Error(err, true, "failed to publish create command")
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	return &gen.PostCreateResponse{GroupID: groupID}, nil
}
