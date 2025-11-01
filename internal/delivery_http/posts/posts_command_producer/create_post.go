package posts_command_producer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandHandler) CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error) {
	// ФЕТЧ БОТА
	// bot, err := p.fetcher.GetBot(ctx, req.BotId.String())
	// if err != nil {
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}

	// ФЕТЧ ПРОФИЛЕЙ (что-то типа)

	// idsString := make([]string, 0, len(req.ProfileIds))
	// for _, id := range req.ProfileIds {
	//	idsString = append(idsString, id.String())
	//}
	//
	// profilesBatch, err := p.fetcher.BatchGetProfiles(ctx, idsString)
	// if err != nil {
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}

	mockedBot := struct {
		Id        uuid.UUID
		BotPrompt string
	}{
		req.BotId,
		"You are a test assistant.",
	}

	mockedProfiles := []struct {
		Id            uuid.UUID
		ProfilePrompt string
	}{
		{
			uuid.New(),
			"Testing. Just say hi and hello world and nothing else.",
		},
		{
			uuid.New(),
			"Add a word test in the end",
		},
	}

	postProfiles := make([]posts.CreatePostProfiles, 0, len(mockedProfiles))
	for _, profile := range mockedProfiles {
		postProfiles = append(postProfiles, posts.CreatePostProfiles{
			ProfileID:     profile.Id,
			ProfilePrompt: profile.ProfilePrompt,
		})
	}

	createPostRequest := posts.KafkaCreatePostRequest{
		PostID:     uuid.New(), // Создаем uuid поста тут
		BotID:      mockedBot.Id,
		BotPrompt:  mockedBot.BotPrompt,
		UserPrompt: req.TaskText,
		Profiles:   postProfiles,
		Platform:   model.PlatformType(req.Platform),
		PostType:   model.PostType(req.PostType.Value),
	}

	rawReq, err := jsoniter.Marshal(createPostRequest)
	if err != nil {
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	// TODO: В рамках теста пока будет создаваться один пост
	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdCreate, createPostRequest.PostID.String(), rawReq), 10*time.Second)
	if err != nil {
		fmt.Println("Ошибка при запросе", err)
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	// TODO: Оставить для timeout-а RAG-a
	//	switch {
	//	case strings.Contains(resp.Error, constants.InternalMsg):
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusNotFound, Message: constants.InternalMsg}, nil
	//}

	if strings.Contains(resp.Error, constants.InternalMsg) {
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusNotFound,
			Message:   constants.InternalMsg,
		}, nil
	}

	returnedPosts := []gen.Post{*resp.PostCommandToGen()}

	return &gen.PostList{
		Data: returnedPosts,
	}, nil
}
