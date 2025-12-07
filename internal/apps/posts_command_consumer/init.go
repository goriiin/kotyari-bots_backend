package posts_command_consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_consumer_client"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts/posts_command_consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/consumer"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka/producer"
	postsRepoLib "github.com/goriiin/kotyari-bots_backend/internal/repo/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/evals"
	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
	"github.com/goriiin/kotyari-bots_backend/pkg/otvet"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"github.com/goriiin/kotyari-bots_backend/pkg/rewriter"
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

func NewPostsCommandConsumer(config *PostsCommandConsumerConfig, llmConfig *LLMConfig) (*PostsCommandConsumer, error) {
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

	cfgRewriter := rewriter.Config{
		NumRewrites: 5,
		Timeout:     60 * time.Second,
	}

	grokClient, err := grok.NewGrokClient(&llmConfig.LLM, &llmConfig.Proxy)
	if err != nil {
		return nil, err
	}

	rw := rewriter.NewGrokRewriter(cfgRewriter, grokClient, "grok-4")

	cfg := evals.Config{
		Timeout: 60 * time.Second,
		Model:   "grok-2-mini",
	}

	j := evals.NewJudge(cfg, grokClient)

	otvetClient, err := otvet.NewOtvetClient(&config.Otvet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create otvet client")
	}

	return &PostsCommandConsumer{
		consumerRunner: posts_command_consumer.NewPostsCommandConsumer(cons, repo, grpc, rw, j, otvetClient),
		consumer:       cons,
		config:         config,
	}, nil
}
