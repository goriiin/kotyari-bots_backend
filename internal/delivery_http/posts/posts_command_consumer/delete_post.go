package posts_command_consumer

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	"github.com/json-iterator/go"
)

func (p *PostsCommandConsumer) DeletePost(ctx context.Context, payload []byte) error {
	fmt.Println("DELETE POST cons")

	var req posts.KafkaDeletePostRequest
	err := jsoniter.Unmarshal(payload, &req)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	err = p.repo.DeletePost(ctx, req.PostID)
	if err != nil {
		return errors.Wrap(err, "failed to delete post")
	}

	return nil
}
