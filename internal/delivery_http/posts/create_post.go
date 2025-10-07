package posts

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"golang.org/x/sync/errgroup"
)

func (p *PostsHandler) CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error) {
	//bot, err := p.fetcher.GetBot(ctx, req.BotId.String())
	//if err != nil {
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}
	mockedBot := struct {
		Id        string
		BotPrompt string
	}{
		uuid.New().String(),
		"You are a test assistant.",
	}

	mockedProfiles := []struct {
		Id            string
		ProfilePrompt string
	}{
		{
			uuid.New().String(),
			"Testing. Just say hi and hello world and nothing else.",
		},
		{
			uuid.New().String(),
			"Add a word test in the end",
		},
	}

	//idsString := make([]string, 0, len(req.ProfileIds))
	//for _, id := range req.ProfileIds {
	//	idsString = append(idsString, id.String())
	//}
	//
	//profilesBatch, err := p.fetcher.BatchGetProfiles(ctx, idsString)
	//if err != nil {
	//	return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}
	//
	//profiles := profilesBatch.GetProfile()

	postsFromInput := make([]model.Post, 0, len(mockedProfiles))
	var postsMutex sync.Mutex

	g, _ := errgroup.WithContext(ctx)
	for _, profile := range mockedProfiles {
		g.Go(func() error {
			post, taskText := httpInputToModel(*req)
			postText, err := p.generator.GeneratePost(ctx, mockedBot.BotPrompt, taskText, profile.ProfilePrompt)
			if err != nil {
				return errors.Wrap(err, "failed to generate post text")
			}
			post.Text = postText
			profileUUID, err := uuid.Parse(profile.Id)
			if err != nil {
				return errors.Wrap(err, "failed to parse profile uuid")
			}
			post.ProfileID = profileUUID
			// TODO: Пока непонятно как отделять заголовок поста от текста - зависит от промта
			post.Title = "XDD"

			postsMutex.Lock()
			postsFromInput = append(postsFromInput, *post)
			postsMutex.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {

		// TODO: наверное, стоит вернуть часть сгенеренных постов
		return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	finalPosts := make([]model.Post, 0, len(postsFromInput))
	for _, inputPost := range postsFromInput {
		post, err := p.repo.CreatePost(ctx, inputPost, req.CategoryIds.Value)
		if err != nil {

			// TODO: должно быть не так явно, но пока сойдет
			return &gen.CreatePostInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
		}
		finalPosts = append(finalPosts, post)
	}

	return &gen.PostList{
		Data: modelSliceToHttpSlice(finalPosts),
	}, nil
}
