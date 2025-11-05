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
	config         *PostsCommandConsumerConfig
}

func NewPostsCommandConsumer(config *PostsCommandConsumerConfig) (*PostsCommandConsumer, error) {
	pool, err := postgres.GetPool(context.Background(), config.Database)
	if err != nil {
		return nil, err
	}

	basicReplier := producer.NewKafkaProducer(&config.KafkaProd)

	cons, err := consumer.NewKafkaRequestReplyConsumer(&config.KafkaCons, basicReplier)
	if err != nil {
		fmt.Println("error happened while creating consumer", err)
		return nil, err
	}

	repo := postsRepoLib.NewPostsCommandRepo(pool)

	grpc, err := posts_consumer_client.NewPostsConsGRPCClient(&config.GRPCServerCfg)
	if err != nil {
		return nil, err
	}

	runner := posts_command_consumer.NewPostsCommandConsumer(cons, repo, grpc)

	return &PostsCommandConsumer{
		consumerRunner: runner,
		consumer:       cons,
		config:         config,
	}, nil
}
