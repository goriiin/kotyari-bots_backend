package posts_command_producer

import (
	"context"
	"fmt"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_producer_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_command_producer"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/producer"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
)

const serviceName = "posts-producer"

type postsCommandHandler interface {
	CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error)
	UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error)
	DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error)
}

type requester interface {
	Request(ctx context.Context, env kafka.Envelope, timeout time.Duration) ([]byte, error)
	Close() error
	Publish(ctx context.Context, env kafka.Envelope) error
}

type PostsCommandProducerApp struct {
	handler  postsCommandHandler
	producer requester
	config   *PostsCommandProducerConfig
	log      *logger.Logger
}

func NewPostsCommandProducerApp(config *PostsCommandProducerConfig) (*PostsCommandProducerApp, error) {
	log := logger.NewLogger(serviceName, &config.ConfigBase)

	grpc, err := posts_producer_client.NewPostsProdGRPCClient(&config.GRPCServerCfg)
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
		handler:  handler,
		producer: p,
		config:   config,
		log:      log,
	}, nil
}
