package posts_command_consumer

import (
	"context"
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_consumer_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_command_consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/producer"
	postsRepoLib "github.com/goriiin/kotyari-bots_backend/internal/repo/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type consumerRunner interface {
	HandleCommands() error
}

type kafkaConsumer interface {
	Start(ctx context.Context) <-chan kafka.CommittableMessage
	Close() error
}

type PostsCommandConsumer struct {
	consumerRunner consumerRunner
	consumer       kafkaConsumer
}

func NewPostsCommandConsumer() (*PostsCommandConsumer, error) {
	pool, err := postgres.GetPool(context.Background(), postgres.Config{
		Host:     "posts_db",
		Port:     5432,
		Name:     "posts",
		User:     "postgres",
		Password: "123",
	})
	if err != nil {
		return nil, err
	}

	basicReplier := producer.NewKafkaProducer(&kafka.KafkaConfig{
		Kind:    "producer",
		Brokers: []string{"kafka:29092"},
		Topic:   "posts-replies",
	})

	cons, err := consumer.NewKafkaRequestReplyConsumer([]string{"kafka:29092"}, "posts-topic", "posts-group-2", basicReplier)
	if err != nil {
		fmt.Println("error happened while creating consumer", err)
		return nil, err
	}

	repo := postsRepoLib.NewPostsCommandRepo(pool)

	grpcClientCfg := &posts_consumer_client.PostsConsGRPCClientConfig{
		ConfigBase: config.ConfigBase{},
		PostsAddr:  "localhost:8080",
		Timeout:    10,
	}
	grpc, _ := posts_consumer_client.NewPostsConsGRPCClient(grpcClientCfg)

	runner := posts_command_consumer.NewPostsCommandConsumer(cons, repo, grpc)

	return &PostsCommandConsumer{
		consumerRunner: runner,
		consumer:       cons,
	}, nil
}
