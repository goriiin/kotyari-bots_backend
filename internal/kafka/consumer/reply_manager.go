package consumer

import (
	"context"
	"fmt"
	"sync"

	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/segmentio/kafka-go"
)

type reader interface {
	GetMessage(ctx context.Context) (kafka.Message, error)
}

type ReplyManager struct {
	mu            sync.RWMutex
	responseChans map[string]chan kafka.Message
	reader        reader
}

func NewReplyManager(reader reader) *ReplyManager {
	return &ReplyManager{
		responseChans: make(map[string]chan kafka.Message),
		reader:        reader,
	}
}

func (rm *ReplyManager) Register(correlationID string) <-chan kafka.Message {
	ch := make(chan kafka.Message, 1)
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.responseChans[correlationID] = ch
	return ch
}

func (rm *ReplyManager) Unregister(correlationID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.responseChans, correlationID)
}

func (rm *ReplyManager) Dispatch(msg kafka.Message) {
	correlationID := kafkaConfig.GetHeader(msg, "correlation_id")
	if correlationID == "" {
		fmt.Println("Error: Reply message is missing correlation_id, skipping.")
		return
	}

	rm.mu.RLock()
	ch, ok := rm.responseChans[correlationID]
	rm.mu.RUnlock()

	if !ok {
		// TODO: Add errors??
		fmt.Println("Warning: Received reply for an unknown or timed-out correlationID:", correlationID)
		return
	}

	select {
	case ch <- msg:
	default:
		// TODO: тоже error
		fmt.Println("Warning: Could not send reply to blocked channel for correlationID:", correlationID)
	}
}

func (rm *ReplyManager) StartConsumingReplies() {
	ctx := context.Background()
	fmt.Println("Reply manager started. Listening for replies...")

	for {
		msg, err := rm.reader.GetMessage(ctx)
		if err != nil {
			fmt.Println("Reply consumer error:", err)
			continue
		}
		rm.Dispatch(msg)
	}
}
