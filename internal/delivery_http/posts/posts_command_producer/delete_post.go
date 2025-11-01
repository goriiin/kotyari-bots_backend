package posts_command_producer

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/json-iterator/go"
)

func (p *PostsCommandHandler) DeletePostById(ctx context.Context, params gen.DeletePostByIdParams) (gen.DeletePostByIdRes, error) {
	req := posts.KafkaDeletePostRequest{PostID: params.PostId}

	rawReq, err := jsoniter.Marshal(req)
	if err != nil {
		return &gen.DeletePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	env := posts.PayloadToEnvelope(posts.CmdDelete, params.PostId.String(), rawReq)
	rawResp, err := p.producer.Request(ctx, env, 5*time.Second)
	if err != nil {
		// TODO: TIMEOUT / PUBLISH ERR
		return &gen.DeletePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.DeletePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	switch {
	case strings.Contains(resp.Error, constants.InternalMsg):
		return &gen.DeletePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	case strings.Contains(resp.Error, constants.NotFoundMsg):
		return &gen.DeletePostByIdNotFound{
			ErrorCode: http.StatusNotFound,
			Message:   "post not found",
		}, nil
	}

	return &gen.NoContent{}, nil
}
