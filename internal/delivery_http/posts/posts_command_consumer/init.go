package posts_command_consumer

import (
	"context"

	"github.com/google/uuid"
	postssgen "github.com/goriiin/kotyari-bots_backend/api/protos/posts/gen"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"google.golang.org/grpc"
)

type postsGetter interface {
	GetPost(ctx context.Context, userPrompt, profilePrompt, botPrompt string, opts ...grpc.CallOption) (*postssgen.GetPostResponse, error)
}

type repo interface {
	CreatePost(ctx context.Context, post model.Post, categoryIDs []uuid.UUID) (model.Post, error)
	UpdatePost(ctx context.Context, post model.Post) (model.Post, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
	CreatePostsBatch(ctx context.Context, posts []model.Post) (err error)
}

type consumer interface {
	Start(ctx context.Context) <-chan kafkaConfig.CommittableMessage
}

type PostsCommandConsumer struct {
	consumer consumer
	repo     repo
	getter   postsGetter
}

func NewPostsCommandConsumer(consumer consumer, repo repo, getter postsGetter) *PostsCommandConsumer {
	return &PostsCommandConsumer{
		consumer: consumer,
		repo:     repo,
		getter:   getter,
	}
}
