package posts_command_consumer

import (
	"context"
	"fmt"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/json-iterator/go"
)

func (p *PostsCommandConsumer) HandleCommands() error {
	ctx := context.Background()
	fmt.Println("HANDLECOMMANDS")
	for message := range p.consumer.Start(ctx) {
		var env kafkaConfig.Envelope
		err := jsoniter.Unmarshal(message.Msg.Value, &env)
		if err != nil {
			// TODO: Наверное стоит убрать вообще
			// log
			fmt.Println("err unmarshal", err)
			_ = message.Ack(ctx)
			// NACK
		}

		switch env.Command {
		case posts.CmdUpdate:
			fmt.Println("UPD POST COMMAND")

			post, err := p.UpdatePost(ctx, env.Payload)

			var packed posts.KafkaResponse
			if err != nil {
				packed = posts.KafkaResponse{Error: err.Error()}

				resp, _ := jsoniter.Marshal(packed)
				fmt.Println("ERR, ACK ", err)
				err = message.Reply(ctx, resp)
				if err != nil {
					fmt.Println("failed to reply to message", err.Error())
				}
			}

			packed = posts.KafkaResponse{
				Post: post,
			}

			resp, _ := jsoniter.Marshal(packed)
			fmt.Println("NOT ERROR, sending post")
			err = message.Reply(ctx, resp)
			if err != nil {
				fmt.Println("failed to reply to message", err.Error())
			}

		case posts.CmdDelete:
			fmt.Println("DEL POST COMMAND")
			err = p.DeletePost(ctx, env.Payload)

			var packed posts.KafkaResponse
			if err != nil {
				packed = posts.KafkaResponse{Error: err.Error()}

				resp, _ := jsoniter.Marshal(packed)
				fmt.Println("ERR, ACK ", err)
				err = message.Reply(ctx, resp)
				if err != nil {
					fmt.Println("failed to reply to message", err.Error())
				}
			}

			resp, _ := jsoniter.Marshal(posts.KafkaResponse{})
			fmt.Println("NOT ERROR, SUCCESSFULLY DELETED")
			err = message.Reply(ctx, resp)
			if err != nil {
				fmt.Println("failed to reply to message", err.Error())
			}

		case posts.CmdCreate:
			fmt.Println("CREATE POST COMMAND")

			post, err := p.CreatePost(ctx, env.Payload)

			// TODO: Handle RAG timeout error

			var packed posts.KafkaResponse
			if err != nil {
				packed = posts.KafkaResponse{Error: err.Error()}

				resp, _ := jsoniter.Marshal(packed)
				fmt.Println("ERR, ACK ", err)
				err = message.Reply(ctx, resp)
				if err != nil {
					fmt.Println("failed to reply to message", err.Error())
				}
			}

			packed = posts.KafkaResponse{
				Post: post,
			}

			resp, _ := jsoniter.Marshal(packed)
			fmt.Println("NOT ERROR, sending post")
			err = message.Reply(ctx, resp)
			if err != nil {
				fmt.Println("failed to reply to message", err.Error())
			}
		}
	}

	return nil
}
