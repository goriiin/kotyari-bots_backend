package posts

import (
	"context"

	"github.com/google/uuid"
	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/profiles/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"google.golang.org/grpc"
)

type postsRepository interface {
	CreatePost(ctx context.Context, post model.Post, categoryIDs []uuid.UUID) (model.Post, error)
	UpdatePost(ctx context.Context, post model.Post) (model.Post, error)
	GetWithCategories(ctx context.Context, id uint64) (model.PostWithCategories, error)
	GetByID(ctx context.Context, id uint64) (model.Post, error)
	DeletePost(ctx context.Context, id uint64) error
}

type botsProfilesFetcher interface {
	GetBot(ctx context.Context, id string, opts ...grpc.CallOption) (*botsgen.Bot, error)
	GetProfile(ctx context.Context, id string, opts ...grpc.CallOption) (*profilesgen.Profile, error)
	BatchGetProfiles(ctx context.Context, ids []string, opts ...grpc.CallOption) (*profilesgen.BatchGetProfilesResponse, error)
}

type postsGenerator interface {
	GeneratePost(ctx context.Context, botPrompt, taskText, profilePrompt string) (string, error)
}

type PostsHandler struct {
	repo      postsRepository
	fetcher   botsProfilesFetcher
	generator postsGenerator
	// log
}

func NewPostsHandler(repo postsRepository, fetcher botsProfilesFetcher, generator postsGenerator) *PostsHandler {
	return &PostsHandler{
		repo:      repo,
		fetcher:   fetcher,
		generator: generator,
	}
}
