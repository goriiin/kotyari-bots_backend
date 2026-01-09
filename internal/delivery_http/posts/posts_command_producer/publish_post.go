package posts_command_producer

import (
	"context"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandHandler) PublishPost(ctx context.Context, req *gen.PublishPostRequest, params gen.PublishPostParams) (gen.PublishPostRes, error) {
	if !req.Approved {
		return &gen.PublishPostBadRequest{
			ErrorCode: http.StatusBadRequest,
			Message:   "post must be approved to publish",
		}, nil
	}

	publishRequest := posts.KafkaPublishPostRequest{
		PostID:   params.PostId,
		Approved: req.Approved,
	}

	rawReq, err := jsoniter.Marshal(publishRequest)
	if err != nil {
		p.log.Error(err, true, "PublishPost: marshal")
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdPublish, params.PostId.String(), rawReq), 5*time.Second)
	if err != nil {
		p.log.Error(err, true, "PublishPost: request")
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		p.log.Error(err, true, "PublishPost: unmarshal response")
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	if resp.Error != "" {
		p.log.Warn("PublishPost: response error", errors.New(resp.Error))
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   resp.Error,
		}, nil
	}

	return &gen.PublishPostResponse{
		Success: true,
		Message: gen.NewOptString("Post approved for publishing"),
	}, nil
}
