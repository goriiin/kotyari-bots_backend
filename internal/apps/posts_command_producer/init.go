package posts_command_producer

import (
	"context"
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/adapters/auth"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_producer_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_command_producer"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/producer"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
)

type postsCommandHandler interface {
	CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error)
	UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error)
	DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error)
	PublishPost(ctx context.Context, req *gen.PublishPostRequest, params gen.PublishPostParams) (gen.PublishPostRes, error)
	SeenPosts(ctx context.Context, req *gen.PostsSeenRequest) (gen.SeenPostsRes, error)
}

type requester interface {
	Request(ctx context.Context, env kafka.Envelope, timeout time.Duration) ([]byte, error)
	Close() error
}

type PostsCommandProducerApp struct {
	handler    postsCommandHandler
	producer   requester
	config     *PostsCommandProducerConfig
	authClient *auth.Client
}

type securityHandler struct {
	authClient *auth.Client
}

func (s *securityHandler) HandleSessionAuth(ctx context.Context, operationName string, t gen.SessionAuth) (context.Context, error) {
	userID, err := s.authClient.VerifySession(ctx, t.APIKey)
	if err != nil {
		return nil, err
	}
	return user.WithID(ctx, userID), nil
}

func NewPostsCommandProducerApp(config *PostsCommandProducerConfig) (*PostsCommandProducerApp, error) {
	grpc, err := posts_producer_client.NewPostsProdGRPCClient(&config.GRPCServerCfg)
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger("posts-command-producer", &config.ConfigBase)

	authClient, err := auth.NewClient(config.Auth, log)
	if err != nil {
		return nil, err
	}

	reader := consumer.NewKafkaConsumer(log, &config.KafkaCons)
	repliesDispatcher := consumer.NewReplyManager(reader)

	p, err := producer.NewKafkaRequestReplyProducer(&config.KafkaProd, &config.KafkaCons, repliesDispatcher)
	if err != nil {
		fmt.Println("error happened while creating producer", err.Error())
		return nil, err
	}

	handler := posts_command_producer.NewPostsHandler(grpc, p, log)

	return &PostsCommandProducerApp{
		handler:    handler,
		producer:   p,
		config:     config,
		authClient: authClient,
	}, nil
}
