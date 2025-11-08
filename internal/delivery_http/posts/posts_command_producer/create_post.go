package posts_command_producer

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandHandler) CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error) {
	//bot, err := p.fetcher.GetBot(ctx, req.BotId.String())
	//if err != nil {
	//	fmt.Println("bots")
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: ierrors.GRPCToDomainError(err).Error()}, nil
	//}
	//
	//idsString := make([]string, 0, len(req.ProfileIds))
	//for _, id := range req.ProfileIds {
	//	idsString = append(idsString, id.String())
	//}
	//
	//profilesBatch, err := p.fetcher.GetProfiles(ctx, idsString)
	//if err != nil {
	//	fmt.Println("profiles")
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: ierrors.GRPCToDomainError(err).Error()}, nil
	//}

	mockedBot := struct {
		Id        uuid.UUID
		BotPrompt string
		BotName   string
	}{
		req.BotId,
		"Промт бота",
		"Крутой бот",
	}

	mockedProfiles := []struct {
		Id            uuid.UUID
		ProfilePrompt string
		ProfileName   string
	}{
		{
			uuid.New(),
			"Крутой промт профиля",
			"Профиль 1",
		},
		{
			uuid.New(),
			"Супер-пупер промт",
			"Профиль 2",
		},
	}

	postProfiles := make([]posts.CreatePostProfiles, 0, len(mockedProfiles))
	for _, profile := range mockedProfiles {
		//profileID, _ := uuid.Parse(profile.Id)
		postProfiles = append(postProfiles, posts.CreatePostProfiles{
			ProfileID:     profile.Id,
			ProfilePrompt: profile.ProfilePrompt,
			ProfileName:   profile.ProfileName,
		})
	}

	groupID := uuid.New()
	//botID, _ := uuid.Parse(bot.Id)
	createPostRequest := posts.KafkaCreatePostRequest{
		GroupID:    groupID,
		BotID:      mockedBot.Id,
		BotName:    mockedBot.BotName,
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

	err = p.producer.Publish(ctx, posts.PayloadToEnvelope(posts.CmdCreate, createPostRequest.GroupID.String(), rawReq))
	if err != nil {
		return &gen.CreatePostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}
	//// TODO: В рамках теста пока будет создаваться один пост
	//rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdCreate, createPostRequest.PostID.String(), rawReq), 10*time.Second)
	//if err != nil {
	//	fmt.Println("Ошибка при запросе", err)
	//	return &gen.CreatePostInternalServerError{
	//		ErrorCode: http.StatusInternalServerError,
	//		Message:   err.Error(),
	//	}, nil
	//}
	//
	//var resp posts.KafkaResponse
	//err = jsoniter.Unmarshal(rawResp, &resp)
	//if err != nil {
	//	return &gen.CreatePostInternalServerError{
	//		ErrorCode: http.StatusInternalServerError,
	//		Message:   err.Error(),
	//	}, nil
	//}

	// TODO: Оставить для timeout-а RAG-a
	//	switch {
	//	case strings.Contains(resp.Error, constants.InternalMsg):
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusNotFound, Message: constants.InternalMsg}, nil
	//}

	//if strings.Contains(resp.Error, constants.InternalMsg) {
	//	return &gen.CreatePostInternalServerError{
	//		ErrorCode: http.StatusInternalServerError,
	//		Message:   constants.InternalMsg,
	//	}, nil
	//}
	//
	//returnedPosts := []gen.Post{*resp.PostCommandToGen()}

	return &gen.PostCreateResponse{GroupID: groupID}, nil
}
