package posts_command_producer

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandHandler) SeenPosts(ctx context.Context, req *gen.PostsSeenRequest) (gen.SeenPostsRes, error) {
	seenPostsRequest := posts.KafkaSeenPostsRequest{
		PostIDs: req.Seen,
	}

	rawReq, err := jsoniter.Marshal(seenPostsRequest)
	if err != nil {
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdSeen, uuid.New().String(), rawReq), 10*time.Second)
	if err != nil {
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}, nil
	}

	switch {
	case strings.Contains(resp.Error, constants.InternalMsg):
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil

	case strings.Contains(resp.Error, constants.NotFoundMsg):
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusNotFound,
			Message:   "post not found",
		}, nil
	}

	return &gen.NoContent{}, nil
}
