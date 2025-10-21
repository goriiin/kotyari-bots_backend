package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type Replier interface {
	Publish(ctx context.Context, message kafka.Message) error
}

type KafkaRequestReplyConsumer struct {
	reader           *kafka.Reader
	config           *kafkaConfig.KafkaConfig
	replier          Replier
	maxCreateRetries int
	baseBackoff      time.Duration
}

// NewKafkaRequestReplyConsumer TODO: Разобраться с инитом с помощью конфига
func NewKafkaRequestReplyConsumer(brokers []string, topic, groupID string, replier Replier) *KafkaRequestReplyConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &KafkaRequestReplyConsumer{
		reader:           r,
		replier:          replier,
		maxCreateRetries: 5,
		baseBackoff:      500 * time.Millisecond,
	}
}

func (c *KafkaRequestReplyConsumer) Start(ctx context.Context) <-chan kafkaConfig.CommittableMessage {
	out := make(chan kafkaConfig.CommittableMessage)
	go func() {
		defer close(out)
		defer c.reader.Close()

		for {
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				// log.Printf("fetch err: %v", err)
				return
			}

			corrID := headerOf(m, "correlation_id")

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
				Reply: func(rctx context.Context, body []byte) error {
					headers := []kafka.Header{
						{Key: "correlation_id", Value: []byte(corrID)},
					}

					return c.replier.Publish(rctx, kafka.Message{
						Value:   body,
						Headers: headers,
					})
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

func headerOf(m kafka.Message, k string) string {
	for _, h := range m.Headers {
		if h.Key == k {
			return string(h.Value)
		}
	}
	return ""
}

//func (c *KafkaRequestReplyConsumer) Start(ctx context.Context) <-chan CommittableMessage {
//	out := make(chan CommittableMessage) // unbuffered = natural backpressure/in-order
//	go func() {
//		defer close(out)
//		defer c.reader.Close()
//
//		for {
//			m, err := c.reader.FetchMessage(ctx) // manual fetch, no auto-commit
//			if err != nil {
//				if errors.Is(err, context.Canceled) {
//					return
//				}
//				// You can choose to log and continue, or terminate:
//				// log.Printf("fetch err: %v", err)
//				return
//			}
//
//			// A one-shot channel so Start waits for your decision (ack/nack)
//			done := make(chan error, 1)
//
//			cm := CommittableMessage{
//				Msg: m,
//				Ack: func(commitCtx context.Context) error {
//					// signal success back to Start loop
//					done <- nil
//					return nil
//				},
//				Nack: func(_ context.Context, _ error) error {
//					// signal failure back to Start loop; we WON'T commit
//					done <- fmt.Errorf("nack")
//					return nil
//				},
//			}
//
//			// Deliver to external handler
//			select {
//			case out <- cm:
//			case <-ctx.Done():
//				return
//			}
//
//			// Wait for external handler's decision
//			select {
//			case decideErr := <-done:
//				if decideErr == nil {
//					// Ack path: commit offset now
//					if err := c.reader.CommitMessages(ctx, m); err != nil {
//						// If commit fails you might want to return (so it can be retried on restart)
//						// or log and continue depending on your tolerance.
//						// log.Printf("commit err: %v", err)
//						return
//					}
//				} else {
//					// Nack path: DO NOT commit. The same message will be re-delivered
//					// after a restart or rebalance. If you prefer "retry topics", push to retry/DLQ here,
//					// then commit to advance (see variant below).
//				}
//			case <-ctx.Done():
//				return
//			}
//		}
//	}()
//	return out
//}
//
//func (c *KafkaRequestReplyConsumer) handleMessage(ctx context.Context, m kafka.Message) error {
//	var env kafkaConfig.Envelope
//	if err := json.Unmarshal(m.Value, &env); err != nil {
//		// TODO: Можно в DLQ, но наверное бессмысленно там держать корраптнутые сообщения
//		return fmt.Errorf("unmarshal: %w", err)
//	}
//
//	var opErr error
//
//	switch env.Command {
//	case kafkaConfig.CmdUpdate:
//		e := decodeUpdate(env.Payload) // your code
//		opErr = c.db.UpdateEntity(ctx, e)
//	case kafkaConfig.CmdDelete:
//		opErr = c.db.DeleteEntity(ctx, env.EntityID)
//	case kafkaConfig.CmdCreate:
//		// idempotency tip: first write a row keyed by CorrelationID (unique),
//		// so duplicates won't double-insert if this handler is retried.
//		opErr = c.handleCreate(ctx, env)
//	default:
//		return fmt.Errorf("unknown command: %s", env.Command)
//	}
//
//	// Send reply if requested (return status to producer)
//	if env.ReplyTo != "" {
//		status := "ok"
//		var msg string
//		if opErr != nil {
//			status, msg = "error", opErr.Error()
//		}
//
//		resp := struct {
//			CorrelationID string `json:"correlation_id"`
//			Status        string `json:"status"`
//			Message       string `json:"message,omitempty"`
//			EntityID      string `json:"entity_id"`
//			Command       string `json:"command"`
//		}{
//			CorrelationID: env.CorrelationID,
//			Status:        status,
//			Message:       msg,
//			EntityID:      env.EntityID,
//			Command:       env.Command,
//		}
//		b, _ := json.Marshal(resp)
//		_ = c.replier.Write(ctx, kafka.Message{
//			Key:   []byte(env.EntityID),
//			Value: b,
//			Headers: []kafka.Header{
//				{Key: "correlation_id", Value: []byte(env.CorrelationID)},
//				{Key: "command", Value: []byte(env.Command)},
//			},
//		})
//	}
//
//	return opErr
//}
//
//func (c *KafkaRequestReplyConsumer) handleCreate(ctx context.Context, env kafkaConfig.Envelope) error {
//	backoff := c.baseBackoff
//
//	for attempt := 0; attempt < c.maxCreateRetries; attempt++ {
//		// 1) external fetch
//		add, err := c.ext.Fetch(ctx, env.Payload)
//		if err == nil {
//			// 2) combine payload + external data and insert (prefer transaction)
//			e := buildEntity(env.Payload, add, env.EntityID)
//			if err = c.db.InsertEntity(ctx, e); err == nil {
//				return nil // success -> caller will commit
//			}
//		}
//
//		// retry path
//		select {
//		case <-time.After(backoff):
//			backoff = minDuration(backoff*2, 10*time.Second)
//			continue
//		case <-ctx.Done():
//			return ctx.Err()
//		}
//	}
//	return fmt.Errorf("create failed after %d retries", c.maxCreateRetries)
//}
//
//func minDuration(a, b time.Duration) time.Duration {
//	if a < b {
//		return a
//	}
//	return b
//}
