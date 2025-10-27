package posts_command_consumer

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/json-iterator/go"
)

func (p *PostsCommandConsumer) UpdatePost(ctx context.Context, payload []byte) (model.Post, error) {
	fmt.Println("UPDATE POST cons")

	var postToUpdate posts.KafkaUpdatePostRequest

	err := jsoniter.Unmarshal(payload, &postToUpdate)
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to unwrap")
	}

	post, err := p.repo.UpdatePost(ctx, postToUpdate.ToModel())
	if err != nil {
		return model.Post{}, errors.Wrap(err, "failed to update post")
	}

	return post, nil
}
