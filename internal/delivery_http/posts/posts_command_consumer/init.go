package posts_command_consumer

import (
	"context"

	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type repo interface {
	UpdatePost(ctx context.Context, post model.Post) (model.Post, error)
}

type consumer interface {
	Start(ctx context.Context) <-chan kafkaConfig.CommittableMessage
}

type PostsCommandConsumer struct {
	consumer consumer
	repo     repo
}

func NewPostsCommandConsumer(consumer consumer, repo repo) *PostsCommandConsumer {
	return &PostsCommandConsumer{
		consumer: consumer,
		repo:     repo,
	}
}
