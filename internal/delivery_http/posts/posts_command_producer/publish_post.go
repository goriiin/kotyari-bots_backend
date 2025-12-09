package posts_command_producer

import (
	"context"
	"net/http"
	"time"

	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
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
		PostID:  params.PostId,
		Approved: req.Approved,
	}

	rawReq, err := jsoniter.Marshal(publishRequest)
	if err != nil {
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdPublish, params.PostId.String(), rawReq), 5*time.Second)
	if err != nil {
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	if resp.Error != "" {
		return &gen.PublishPostInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   resp.Error,
		}, nil
	}

	return &gen.PublishPostResponse{
		Success: true,
		Message: "Post approved for publishing",
	}, nil
}

