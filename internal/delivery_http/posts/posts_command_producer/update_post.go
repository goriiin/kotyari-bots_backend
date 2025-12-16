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

func (p *PostsCommandHandler) UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error) {
	updatePostRequest := posts.KafkaUpdatePostRequest{
		PostID: params.PostId,
		Title:  req.Title,
		Text:   req.Text,
	}

	rawReq, err := jsoniter.Marshal(updatePostRequest)
	if err != nil {
		p.log.Error(err, true, "failed to marshal update request")
		return &gen.UpdatePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdUpdate, params.PostId.String(), rawReq), 5*time.Second)
	if err != nil {
		p.log.Error(err, true, "failed to request update post via kafka")
		return &gen.UpdatePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		p.log.Error(err, true, "failed to unmarshal kafka response")
		return &gen.UpdatePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	switch {
	case strings.Contains(resp.Error, constants.InternalMsg):
		p.log.Error(nil, false, "received internal error from kafka reply: "+resp.Error)
		return &gen.UpdatePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil

	case strings.Contains(resp.Error, constants.NotFoundMsg):
		return &gen.UpdatePostByIdNotFound{
			ErrorCode: http.StatusNotFound,
			Message:   "post not found",
		}, nil
	}

	return resp.PostCommandToGen(), nil
}
