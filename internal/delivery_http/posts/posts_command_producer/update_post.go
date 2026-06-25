package posts_command_producer

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	"github.com/json-iterator/go"
)

func (p *PostsCommandHandler) UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return nil, err
	}

	updatePostRequest := posts.KafkaUpdatePostRequest{
		PostID: params.PostId,
		UserID: userID,
		Title:  req.Title,
		Text:   req.Text,
	}

	rawReq, err := jsoniter.Marshal(updatePostRequest)
	if err != nil {
		p.log.Error(err, true, "UpdatePostById: marshal")
		return &gen.UpdatePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdUpdate, params.PostId.String(), rawReq), 5*time.Second)
	if err != nil {
		p.log.Error(err, true, "UpdatePostById: request")
		return &gen.UpdatePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		p.log.Error(err, true, "UpdatePostById: unmarshal response")
		return &gen.UpdatePostByIdInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	switch {
	case strings.Contains(resp.Error, constants.InternalMsg):
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
