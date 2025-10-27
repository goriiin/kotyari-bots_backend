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

type reader interface {
	GetMessage(ctx context.Context) (kafka.Message, error)
}

type KafkaRequestReplyProducer struct {
	writer     *kafka.Writer
	config     *kafkaConfig.KafkaConfig
	reader     reader
	replyTopic string
	replyGroup string
}

func NewKafkaRequestReplyProducer(config *kafkaConfig.KafkaConfig, replyTopic, replyGroup string, reader reader) (*KafkaRequestReplyProducer, error) {
	// TODO: XDDDDDDD
	if err := kafkaConfig.EnsureTopicCreated(config.Brokers[0], replyTopic); err != nil {
		fmt.Println("ТОПИК НЕ СОЗДАЛСЯ))) РАЗРАБ КАФКА-ГО МОЛОДЕЦ")
	}

	return &KafkaRequestReplyProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(config.Brokers...),
			Topic:                  config.Topic,
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
		},
		replyTopic: replyTopic,
		replyGroup: replyGroup,
		config:     config,
		reader:     reader,
	}, nil
}

func (p *KafkaRequestReplyProducer) Publish(ctx context.Context, env *kafkaConfig.Envelope) error {
	env.CorrelationID = uuid.NewString()
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
	fmt.Println("Зашли в publish: ", time.Now())
	if err := p.Publish(ctx, &env); err != nil {
		return nil, err
	}
	fmt.Println("Вышли из publish: ", time.Now())

	fmt.Printf("CFG PROD: %+v\n", p.config)
	fmt.Printf("CFG READ: repl topic: %v, replc group: %v", p.replyTopic, p.replyGroup)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	fmt.Println("COR ID", env.CorrelationID)

	for {
		// TODO: Надо поменять на Fetch + Commit потому что хуесос я
		m, err := p.reader.GetMessage(ctx)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if getHeader(m, "correlation_id") == env.CorrelationID {
			fmt.Printf("MESSAGE RECEIVED: %+v\n", m)
			fmt.Printf("CONTEXT VNUTRI MAMASHI: %+v\n", ctx)
			fmt.Println("Получили ответ: ", time.Now())
			return m.Value, nil
		}
	}
}

func getHeader(m kafka.Message, key string) string {
	for _, h := range m.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}
