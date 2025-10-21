package posts_command_consumer

import (
	"context"
	"encoding/json"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
)

func (p *PostsCommandConsumer) UpdatePost(ctx context.Context, payload []byte) error {
	var postToUpdate posts.RawPostUpdate

	err := json.Unmarshal(payload, &postToUpdate)
	if err != nil {
		errors.Wrap(err, "failed to unwrap")
	}

	_, err = p.repo.UpdatePost(ctx, posts.FromRawPost(postToUpdate))
	if err != nil {
		return errors.Wrap(err, "failed to update post")
	}

	return nil
}
