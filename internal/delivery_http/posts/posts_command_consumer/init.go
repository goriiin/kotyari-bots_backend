package posts_command_consumer

import (
	"context"

	"github.com/google/uuid"
	postssgen "github.com/goriiin/kotyari-bots_backend/api/protos/posts/gen"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/otvet"
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

type otvetClient interface {
	CreatePost(ctx context.Context, req *otvet.CreatePostRequest) (*otvet.CreatePostResponse, error)
	CreatePostSimple(ctx context.Context, title string, contentText string, topicType int, spaces []otvet.Space) (*otvet.CreatePostResponse, error)
	PredictTagsSpaces(ctx context.Context, text string) (*otvet.PredictTagsSpacesResponse, error)
}

type PostsCommandConsumer struct {
	consumer    consumer
	repo        repo
	getter      postsGetter
	rewriter    rewriter
	judge       judge
	otvetClient otvetClient
}

func NewPostsCommandConsumer(
	consumer consumer,
	repo repo,
	getter postsGetter,
	rewriter rewriter,
	judge judge,
	otvetClient otvetClient,
) *PostsCommandConsumer {
	return &PostsCommandConsumer{
		consumer:    consumer,
		repo:        repo,
		getter:      getter,
		rewriter:    rewriter,
		judge:       judge,
		otvetClient: otvetClient,
	}
}
