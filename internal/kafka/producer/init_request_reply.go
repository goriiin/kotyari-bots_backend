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
	StartConsumingReplies()
	Dispatch(msg kafka.Message)
	Register(correlationID string) <-chan kafka.Message
	Unregister(correlationID string)
}

type KafkaRequestReplyProducer struct {
	writer     *kafka.Writer
	config     *kafkaConfig.KafkaConfig
	replyTopic string
	replyGroup string
	dispatcher repliesDispatcher
}

func NewKafkaRequestReplyProducer(config *kafkaConfig.KafkaConfig, replyTopic, replyGroup string, dispatcher repliesDispatcher) (*KafkaRequestReplyProducer, error) {
	// TODO: XDDDDDDD
	if err := kafkaConfig.EnsureTopicCreated(config.Brokers[0], replyTopic); err != nil {
		fmt.Println("Failed to create topic", err.Error())
	}

	producer := &KafkaRequestReplyProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(config.Brokers...),
			Topic:                  config.Topic,
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
		},
		replyTopic: replyTopic,
		replyGroup: replyGroup,
		config:     config,
		dispatcher: dispatcher,
	}

	go producer.dispatcher.StartConsumingReplies()
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
			{Key: "reply_to", Value: []byte(p.replyTopic)},
		},
	})
}

func (p *KafkaRequestReplyProducer) Request(ctx context.Context, env kafkaConfig.Envelope, timeout time.Duration) ([]byte, error) {
	env.CorrelationID = uuid.NewString()
	replyChan := p.dispatcher.Register(env.CorrelationID)
	defer p.dispatcher.Unregister(env.CorrelationID)

	fmt.Println("Зашли в publish: ", time.Now())
	if err := p.Publish(ctx, env); err != nil {
		return nil, err
	}
	fmt.Println("Вышли из publish: ", time.Now())

	fmt.Printf("CFG PROD: %+v\n", p.config)
	fmt.Printf("CFG READ: repl topic: %v, replc group: %v", p.replyTopic, p.replyGroup)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	fmt.Println("COR ID", env.CorrelationID)

	select {
	case msg := <-replyChan:
		fmt.Printf("Received reply for CorrelationID: %s\n", env.CorrelationID)
		return msg.Value, nil
	case <-ctx.Done():
		fmt.Printf("Timed out waiting for reply for CorrelationID: %s\n", env.CorrelationID)
		return nil, ctx.Err()
	}
}
