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
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	postsRepoLib "github.com/goriiin/kotyari-bots_backend/internal/repo/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/evals"
	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
	"github.com/goriiin/kotyari-bots_backend/pkg/otvet"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"github.com/goriiin/kotyari-bots_backend/pkg/posting_queue"
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
	log := logger.NewLogger("posts-command-consumer", &config.ConfigBase)

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

	// Initialize posting queue
	postingInterval := config.PostingQueue.PostingInterval
	if postingInterval == 0 {
		postingInterval = 30 * time.Minute
	}

	processingInterval := config.PostingQueue.ProcessingInterval
	if processingInterval == 0 {
		processingInterval = 1 * time.Minute
	}

	queue := posting_queue.NewQueue(
		postingInterval,
		processingInterval,
		config.PostingQueue.ModerationRequired,
	)

	accountID := "default"
	queue.AddAccount(accountID, config.Otvet.AuthToken, otvetClient)

	publishFunc := func(ctx context.Context, account *posting_queue.Account, queuedPost *posting_queue.QueuedPost) error {
		if account.Client == nil {
			log.Error(nil, false, "Queue: account client is nil")
			return errors.New("account client is nil")
		}

		log.Info(fmt.Sprintf("Queue: publishing post %s...", queuedPost.Post.ID))

		otvetResp, err := account.Client.CreatePostSimple(
			ctx,
			queuedPost.Candidate.Title,
			queuedPost.Candidate.Text,
			queuedPost.Request.TopicType,
			queuedPost.Request.Spaces,
		)
		if err != nil {
			log.Error(err, true, fmt.Sprintf("Queue: failed to publish post %s", queuedPost.Post.ID))
			return errors.Wrap(err, "failed to publish post from queue")
		}

		log.Info(fmt.Sprintf("Queue: successfully published post %s. Resp ID: %d", queuedPost.Post.ID, otvetResp.Result.ID))

		if otvetResp != nil && otvetResp.Result != nil {
			// Обновляем модель
			queuedPost.Post.OtvetiID = uint64(otvetResp.Result.ID)
			queuedPost.Post.IsPublished = true
			queuedPost.Post.URL = fmt.Sprintf("https://otvet.mail.ru/question/%d", otvetResp.Result.ID)

			// Сохраняем обновление в БД
			if _, err := repo.UpdatePost(ctx, *queuedPost.Post); err != nil {
				log.Error(err, true, "Queue: failed to update post status in DB after publish")
			}
		}

		return nil
	}

	ctx := context.Background()
	go queue.StartProcessing(ctx, publishFunc)

	return &PostsCommandConsumer{
		consumerRunner: posts_command_consumer.NewPostsCommandConsumer(cons, repo, grpc, rw, j, otvetClient, queue, log),
		consumer:       cons,
		config:         config,
	}, nil
}
