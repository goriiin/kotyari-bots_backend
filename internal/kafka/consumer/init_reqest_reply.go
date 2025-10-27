package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type replier interface {
	Publish(ctx context.Context, message kafka.Message) error
}

type KafkaRequestReplyConsumer struct {
	reader           *kafka.Reader
	config           *kafkaConfig.KafkaConfig
	replier          replier
	maxCreateRetries int
	baseBackoff      time.Duration
}

// NewKafkaRequestReplyConsumer TODO: Разобраться с инитом с помощью конфига
func NewKafkaRequestReplyConsumer(brokers []string, topic, groupID string, replier replier) (*KafkaRequestReplyConsumer, error) {
	if err := kafkaConfig.EnsureTopicCreated(brokers[0], topic); err != nil {
		fmt.Println("ТОПИК СОСАЛ ХУЙ СТРОГО")

		return nil, err
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               brokers,
		Topic:                 topic,
		GroupID:               groupID,
		GroupTopics:           []string{topic},
		MinBytes:              1,
		MaxBytes:              10e6,
		WatchPartitionChanges: true, // ???
	})
	return &KafkaRequestReplyConsumer{
		reader:           r,
		replier:          replier,
		maxCreateRetries: 5,
		baseBackoff:      500 * time.Millisecond,
	}, nil
}

func (c *KafkaRequestReplyConsumer) Start(ctx context.Context) <-chan kafkaConfig.CommittableMessage {
	fmt.Println("START")
	out := make(chan kafkaConfig.CommittableMessage)
	go func() {
		defer close(out)
		defer c.reader.Close()

		for {
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				fmt.Println("reading error", err)
				if errors.Is(err, context.Canceled) {
					return
				}
				// log.Printf("fetch err: %v", err)
				return
			}

			fmt.Printf("ALO ALO %+v\n", m)

			corrID := kafkaConfig.GetHeader(m, "correlation_id")

			fmt.Println("corrID: ", corrID)
			done := make(chan error, 1)

			cm := kafkaConfig.CommittableMessage{
				Msg: m,
				Ack: func(commitCtx context.Context) error {
					done <- nil
					return nil
				},
				Nack: func(_ context.Context, _ error) error {
					done <- fmt.Errorf("nack")
					return nil
				},
				// TODO: Скорее все помимо reply будет еще replyWithError, чтобы не двигать оффсет
				Reply: func(ctx context.Context, body []byte) error {
					done <- nil // оффсет двигается
					headers := []kafka.Header{
						{Key: "correlation_id", Value: []byte(corrID)},
					}

					// TODO: просто ужас (даже не тихий)
					err = c.replier.Publish(ctx, kafka.Message{
						Key:     []byte(corrID),
						Value:   body,
						Headers: headers,
					})
					if err != nil {
						fmt.Println("ERROR PUBLISHING MESSAGE ", err.Error())
					}
					//return c.replier.Publish(ctx, kafka.Message{
					//	Value:   body,
					//	Headers: headers,
					//})
					return nil
				},

				//ReplyWithError: func(ctx context.Context, body []byte) error {
				//
				//},
			}

			select {
			case out <- cm:
			case <-ctx.Done():
				return
			}

			select {
			case decideErr := <-done:
				if decideErr == nil {
					fmt.Println("КОММИТ СООБЩЕНИЙ")

					if err := c.reader.CommitMessages(ctx, m); err != nil {
						// log.Printf("commit err: %v", err)
						return
					}
				} else {
					// Nack path: DO NOT commit. The same message will be re-delivered
					// after a restart or rebalance. If you prefer "retry topics", push to retry/DLQ here,
					// then commit to advance (see variant below).
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}
