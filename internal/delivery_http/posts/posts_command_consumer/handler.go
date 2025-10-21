package posts_command_consumer

import (
	"context"
	"encoding/json"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	kafkaConfig "github.com/goriiin/kotyari-bots_backend/internal/kafka"
)

func (p *PostsCommandConsumer) HandleCommands() error {
	ctx := context.Background() // TODO: Заменить на WithTimeout

	for message := range p.consumer.Start(ctx) {
		var env kafkaConfig.Envelope
		err := json.Unmarshal(message.Msg.Value, &env)
		if err != nil {
			// log
			_ = message.Nack(ctx, err)

		}

		switch env.Command {
		case posts.CmdUpdate:
			err = p.UpdatePost(ctx, env.Payload)
			if err != nil {
				_ = message.Ack(ctx)
			}
		}
	}

	return nil
}
