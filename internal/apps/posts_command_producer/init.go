package posts_command_producer

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_command_producer"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/producer"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type postsCommandHandler interface {
	CreatePost(ctx context.Context, req *gen.PostInput) (gen.CreatePostRes, error)
	CreatePostSEO(ctx context.Context, req *gen.PostInput) (gen.CreatePostSEORes, error)
	UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error)
	DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error)
}

type PostsCommandProducerApp struct {
	handler postsCommandHandler
}

func NewPostsCommandProducerApp() *PostsCommandProducerApp {
	// PIVO
	grpcClientCfg := &posts_client.PostsGRPCClientAppConfig{
		ConfigBase:   config.ConfigBase{},
		BotsAddr:     "localhost:8080",
		ProfilesAddr: "localhost:8081",
		PostsAddr:    "localhost:8082",
		Timeout:      10,
	}

	grpc, _ := posts_client.NewPostsGRPCClient(grpcClientCfg)

	kafkaCfg := &kafka.KafkaConfig{
		Kind:    "producer",
		Brokers: []string{"kafka:29092"},
		Topic:   "posts-topic",
		GroupID: "posts-group",
	}

	p := producer.NewKafkaRequestReplyProducer(kafkaCfg, "posts-replies", "posts-replies-group")

	handler := posts_command_producer.NewPostsHandler(grpc, p)

	return &PostsCommandProducerApp{
		handler: handler,
	}
}
