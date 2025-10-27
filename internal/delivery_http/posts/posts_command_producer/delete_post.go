package posts_command_producer

import (
	"context"
	"net/http"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/json-iterator/go"
)

func (p *PostsCommandHandler) DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error) {
	req := posts.KafkaDeletePostRequest{PostID: params.PostId}

	rawReq, err := jsoniter.Marshal(req)
	if err != nil {
		return &gen.DeletePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	env := posts.PayloadToEnvelope(posts.CmdDelete, params.PostId.String(), rawReq)
	rawResp, err := p.producer.Request(ctx, env, 5*time.Second)
	if err != nil {
		return &gen.DeletePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.DeletePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	if resp.IsError {
		// TODO: Ошибка на стороне consumer, в идеале потом сделать switch через errors.Is и отдельно делать сообщения для каждого случая
		return &gen.DeletePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: resp.StatusMessage}, nil
	}

	return &gen.NoContent{}, nil
}
