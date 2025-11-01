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
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type postsCommandHandler interface {
	CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error)
	UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error)
	DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error)
}

type requester interface {
	Request(ctx context.Context, env kafka.Envelope, timeout time.Duration) ([]byte, error)
	Close() error
}

type PostsCommandProducerApp struct {
	handler  postsCommandHandler
	producer requester
}

func NewPostsCommandProducerApp() (*PostsCommandProducerApp, error) {
	// PIVO
	grpcClientCfg := &posts_producer_client.PostsProdGRPCClientConfig{
		ConfigBase:   config.ConfigBase{},
		BotsAddr:     "localhost:8080",
		ProfilesAddr: "localhost:8081",
		Timeout:      10,
	}

	grpc, _ := posts_producer_client.NewPostsProdGRPCClient(grpcClientCfg)

	kafkaCfg := &kafka.KafkaConfig{
		Kind:    "producer",
		Brokers: []string{"kafka:29092"},
		Topic:   "posts-topic",
		GroupID: "posts-group",
	}

	readerCfg := &kafka.KafkaConfig{
		Kind:    "consumer",
		Brokers: []string{"kafka:29092"},
		Topic:   "posts-replies",
		GroupID: "posts-replies-group",
	}

	log := logger.NewLogger("xdd", &grpcClientCfg.ConfigBase)

	reader := consumer.NewKafkaConsumer(log, readerCfg)

	repliesDispatcher := consumer.NewReplyManager(reader)

	p, err := producer.NewKafkaRequestReplyProducer(kafkaCfg, "posts-replies", "posts-replies-group", repliesDispatcher)
	if err != nil {
		fmt.Println("error happened while creating producer", err.Error())
		return nil, err
	}

	handler := posts_command_producer.NewPostsHandler(grpc, p)

	return &PostsCommandProducerApp{
		handler:  handler,
		producer: p,
	}, nil
}
