package consumer

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-faster/errors"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

const maxBackoff = 10 * time.Second

type replier interface {
	Publish(ctx context.Context, message kafka.Message) error
	Close() error
}

type KafkaRequestReplyConsumer struct {
	reader           *kafka.Reader
	config           *kafkaConfig.KafkaConfig
	replier          replier
	maxCreateRetries int
	baseBackoff      time.Duration
}

func NewKafkaRequestReplyConsumer(config *kafkaConfig.KafkaConfig, replier replier) (*KafkaRequestReplyConsumer, error) {
	// TODO: Check if needed
	if err := kafkaConfig.EnsureTopicCreated(config.Brokers[0], config.Topic); err != nil {
		fmt.Println("failed to create topic")

		return nil, err
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               config.Brokers,
		Topic:                 config.Topic,
		GroupID:               config.GroupID,
		GroupTopics:           []string{config.Topic},
		MinBytes:              1,
		MaxBytes:              10e6,
		WatchPartitionChanges: true, // ???
	})
	return &KafkaRequestReplyConsumer{
		reader:           r,
		replier:          replier,
		maxCreateRetries: 5,
		baseBackoff:      500 * time.Millisecond,
		config:           config,
	}, nil
}

func (c *KafkaRequestReplyConsumer) Start(ctx context.Context) <-chan kafkaConfig.CommittableMessage {
	out := make(chan kafkaConfig.CommittableMessage)

	go func() {
		defer close(out)
		defer func() {
			if err := c.reader.Close(); err != nil {
				fmt.Println("Error closing reader:", err)
			}
		}()

		for {
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				fmt.Println("reading error:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			retryAttempt := 0
			messageProcessed := false

			for !messageProcessed {
				done := make(chan error, 1)
				corrID := kafkaConfig.GetHeader(m, "correlation_id")

				cm := kafkaConfig.CommittableMessage{
					Msg: m,

					Ack: func(commitCtx context.Context) error {
						done <- nil
						return nil
					},

					Nack: func(_ context.Context, err error) error {
						done <- err
						return nil
					},

					Reply: func(replyCtx context.Context, body []byte) error {
						headers := []kafka.Header{{Key: "correlation_id", Value: []byte(corrID)}}
						err := c.replier.Publish(replyCtx, kafka.Message{
							Key:     []byte(corrID),
							Value:   body,
							Headers: headers,
						})
						if err != nil {
							done <- fmt.Errorf("failed to reply: %w", err)
							return err
						}
						done <- nil
						return nil
					},

					ReplyWithError: func(replyCtx context.Context, body []byte) error {
						headers := []kafka.Header{
							{Key: "correlation_id", Value: []byte(corrID)},
						}

						_ = c.replier.Publish(replyCtx, kafka.Message{
							Key:     []byte(corrID),
							Value:   body,
							Headers: headers,
						})

						done <- nil
						return nil
					},
				}

				select {
				case out <- cm:
				case <-ctx.Done():
					return
				}

				select {
				case decideErr := <-done:
					if decideErr == nil {
						if err := c.reader.CommitMessages(ctx, m); err != nil {
							fmt.Println("commit error:", err)
						}
						messageProcessed = true
					} else {
						retryAttempt++
						backoff := c.calculateBackoff(retryAttempt)

						fmt.Printf("Message processing failed (attempt %d). Retrying in %v. Error: %v\n", retryAttempt, backoff, decideErr)

						select {
						case <-time.After(backoff):
						case <-ctx.Done():
							return
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

func (c *KafkaRequestReplyConsumer) Close() error {
	return errors.Join(c.reader.Close(), c.replier.Close())
}

func (c *KafkaRequestReplyConsumer) calculateBackoff(attempt int) time.Duration {
	backoff := c.baseBackoff * time.Duration(math.Pow(2, float64(attempt-1)))
	if backoff > maxBackoff {
		return maxBackoff
	}
	return backoff
}
