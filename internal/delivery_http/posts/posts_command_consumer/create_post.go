package posts_command_consumer

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandConsumer) CreatePost(ctx context.Context, payload []byte) (model.Post, error) {
	var req posts.KafkaCreatePostRequest
	err := jsoniter.Unmarshal(payload, &req)
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to unmarshal")
	}

	// TODO: тут должен быть либо поход батчами c несколькими профилями для создания постов,
	// либо асинк GetPost c wg например, для тестов путь пока будет один
	// post, err := p.getter.GetPost(ctx, req.UserPrompt, req.Profiles[0].ProfilePrompt, req.BotPrompt)

	finalPost := model.Post{
		ID:        req.PostID,
		OtvetiID:  100, // че с этим делать пока непонятно
		BotID:     req.BotID,
		ProfileID: req.Profiles[0].ProfileID, // опять батчи или что-то такое
		Platform:  req.Platform,
		Type:      req.PostType,
		Title:     "Title поста, полученный от геттера",
		Text:      "Text поста, полученный от геттера",
	}

	post, err := p.repo.CreatePost(ctx, finalPost, nil) // Пока без категорий =(
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to create post")
	}

	return post, nil
}
