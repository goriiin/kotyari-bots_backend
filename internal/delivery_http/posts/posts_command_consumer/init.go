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
	UpdatePostsBatch(ctx context.Context, posts []model.Post) (err error)
	SeenPostsBatch(ctx context.Context, postsIds []uuid.UUID) (err error)
}

type consumer interface {
	Start(ctx context.Context) <-chan kafkaConfig.CommittableMessage
}

type rewriter interface {
	Rewrite(ctx context.Context, user, profile, bot string) ([]string, error)
}

type judge interface {
	SelectBest(ctx context.Context, userPrompt, profilePrompt, botPrompt string, candidates []model.Candidate) (model.Candidate, error)
}

type PostsCommandConsumer struct {
	consumer consumer
	repo     repo
	getter   postsGetter
	rewriter rewriter
	judge    judge
}

func NewPostsCommandConsumer(
	consumer consumer,
	repo repo,
	getter postsGetter,
	rewriter rewriter,
	judge judge,
) *PostsCommandConsumer {
	return &PostsCommandConsumer{
		consumer: consumer,
		repo:     repo,
		getter:   getter,
		rewriter: rewriter,
		judge:    judge,
	}
}
