package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

type repliesDispatcher interface {
	StartConsumingReplies(ctx context.Context)
	Dispatch(msg kafka.Message)
	Register(correlationID string) <-chan kafka.Message
	Unregister(correlationID string)
}

type KafkaRequestReplyProducer struct {
	writer     *kafka.Writer
	prodConfig *kafkaConfig.KafkaConfig
	consConfig *kafkaConfig.KafkaConfig
	dispatcher repliesDispatcher
	shutdown   context.CancelFunc
}

func NewKafkaRequestReplyProducer(prodConfig *kafkaConfig.KafkaConfig, consConfig *kafkaConfig.KafkaConfig, dispatcher repliesDispatcher) (*KafkaRequestReplyProducer, error) {
	if err := kafkaConfig.EnsureTopicCreated(prodConfig.Brokers[0], consConfig.Topic); err != nil {
		fmt.Println("Failed to create topic", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	producer := &KafkaRequestReplyProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(prodConfig.Brokers...),
			Topic:                  prodConfig.Topic,
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
		},
		consConfig: consConfig,
		prodConfig: prodConfig,
		dispatcher: dispatcher,
		shutdown:   cancel,
	}

	go producer.dispatcher.StartConsumingReplies(ctx)
	return producer, nil
}

func (p *KafkaRequestReplyProducer) Publish(ctx context.Context, env kafkaConfig.Envelope) error {
	b, err := jsoniter.Marshal(env)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(env.EntityID),
		Value: b,
		Headers: []kafka.Header{
			{Key: "correlation_id", Value: []byte(env.CorrelationID)},
			{Key: "command", Value: []byte(env.Command)},
			{Key: "reply_to", Value: []byte(p.consConfig.Topic)},
		},
	})
}

func (p *KafkaRequestReplyProducer) Request(ctx context.Context, env kafkaConfig.Envelope, timeout time.Duration) ([]byte, error) {
	env.CorrelationID = uuid.NewString()
	replyChan := p.dispatcher.Register(env.CorrelationID)
	defer p.dispatcher.Unregister(env.CorrelationID)

	if err := p.Publish(ctx, env); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case msg := <-replyChan:
		fmt.Printf("Received reply for CorrelationID: %s\n", env.CorrelationID)
		return msg.Value, nil
	case <-ctx.Done():
		fmt.Printf("Timed out waiting for reply for CorrelationID: %s\n", env.CorrelationID)
		return nil, ctx.Err()
	}
}

func (p *KafkaRequestReplyProducer) Close() error {
	p.shutdown()

	return p.writer.Close()
}
