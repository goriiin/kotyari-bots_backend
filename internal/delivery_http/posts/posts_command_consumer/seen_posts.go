package posts_command_consumer

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandConsumer) SeenPosts(ctx context.Context, payload []byte) error {
	var seenPosts posts.KafkaSeenPostsRequest

	err := jsoniter.Unmarshal(payload, &seenPosts)
	if err != nil {
		return errors.Wrap(err, "failed to unwrap")
	}

	err = p.repo.SeenPostsBatch(ctx, seenPosts.PostIDs)
	if err != nil {
		return errors.Wrap(err, "failed to change posts status")
	}

	return nil
}
