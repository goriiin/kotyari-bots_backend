package posts_command_consumer

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/json-iterator/go"
)

const failedToSendReplyMsg = "failed to send reply successfully"

func (p *PostsCommandConsumer) HandleCommands() error {
	ctx := context.Background()
	for message := range p.consumer.Start(ctx) {
		var env kafkaConfig.Envelope
		if err := jsoniter.Unmarshal(message.Msg.Value, &env); err != nil {
			fmt.Printf("%s: %v\n", constants.ErrUnmarshal, err)
			continue
		}

		var err error
		switch env.Command {
		case posts.CmdUpdate:
			err = p.handleUpdateCommand(ctx, message, env.Payload)
		case posts.CmdDelete:
			err = p.handleDeleteCommand(ctx, message, env.Payload)
		case posts.CmdCreate:
			err = p.handleCreateCommand(ctx, message, env.Payload)
		default:
			err = errors.Errorf("unknown command received: %s", env.Command)
		}

		if err != nil {
			fmt.Printf("failed to handle command '%s': %v\n", env.Command, err)
		}
	}

	return nil
}

func (p *PostsCommandConsumer) handleUpdateCommand(ctx context.Context, message kafkaConfig.CommittableMessage, payload []byte) error {
	post, err := p.UpdatePost(ctx, payload)
	if err != nil {
		return sendErrReply(ctx, message, err)
	}

	resp := posts.KafkaResponse{
		Post: post,
	}

	rawResp, err := jsoniter.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, constants.MarshalMsg)
	}

	if err := message.Reply(ctx, rawResp); err != nil {
		return errors.Wrap(err, failedToSendReplyMsg)
	}

	return nil
}

func (p *PostsCommandConsumer) handleDeleteCommand(ctx context.Context, message kafkaConfig.CommittableMessage, payload []byte) error {
	if err := p.DeletePost(ctx, payload); err != nil {
		return sendErrReply(ctx, message, err)
	}

	resp, err := jsoniter.Marshal(posts.KafkaResponse{})
	if err != nil {
		return errors.Wrap(err, constants.MarshalMsg)
	}

	if err := message.Reply(ctx, resp); err != nil {
		return errors.Wrap(err, failedToSendReplyMsg)
	}

	return nil
}

func (p *PostsCommandConsumer) handleCreateCommand(ctx context.Context, message kafkaConfig.CommittableMessage, payload []byte) error {
	err := p.CreatePost(ctx, payload)
	if err != nil {
		// TODO: Handle RAG timeout error specifically if needed, otherwise send generic error.
		return sendErrReply(ctx, message, err)
	}

	if err = message.Ack(ctx); err != nil {
		return errors.Wrap(err, "failed to ACK posts creation")
	}

	return nil
}

func sendErrReply(ctx context.Context, message kafkaConfig.CommittableMessage, originalErr error) error {
	errMessage := posts.KafkaResponse{Error: originalErr.Error()}
	resp, err := jsoniter.Marshal(errMessage)
	if err != nil {
		return errors.Wrap(err, constants.MarshalMsg)
	}

	if err := message.Reply(ctx, resp); err != nil {
		return errors.Wrap(err, "failed to send error reply")
	}

	return nil
}
