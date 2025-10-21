package posts_command_consumer

import (
	"context"

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

type PostsCommandConsumer struct {
	consumerRunner consumerRunner
}

func NewPostsCommandConsumer() (*PostsCommandConsumer, error) {
	pool, err := postgres.GetPool(context.Background(), postgres.Config{
		Host:     "posts_db",
		Port:     54327,
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
		GroupID: "posts-replies-group",
	})

	cons := consumer.NewKafkaRequestReplyConsumer([]string{"kafka:29092"}, "posts-topic", "posts-group", basicReplier)
	repo := postsRepoLib.NewPostsCommandRepo(pool)

	runner := posts_command_consumer.NewPostsCommandConsumer(cons, repo)

	return &PostsCommandConsumer{
		consumerRunner: runner,
	}, nil
}
