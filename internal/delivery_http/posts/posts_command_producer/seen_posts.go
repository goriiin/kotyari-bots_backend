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
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
	jsoniter "github.com/json-iterator/go"
)

func (p *PostsCommandHandler) SeenPosts(ctx context.Context, req *gen.PostsSeenRequest) (gen.SeenPostsRes, error) {
	userID, err := user.GetID(ctx)
	if err != nil {
		return nil, err
	}

	seenPostsRequest := posts.KafkaSeenPostsRequest{
		PostIDs: req.Seen,
		UserID:  userID,
	}

	rawReq, err := jsoniter.Marshal(seenPostsRequest)
	if err != nil {
		p.log.Error(err, true, "SeenPosts: marshal")
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdSeen, uuid.New().String(), rawReq), 10*time.Second)
	if err != nil {
		p.log.Error(err, true, "SeenPosts: request")
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		p.log.Error(err, true, "SeenPosts: unmarshal response")
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil
	}

	switch {
	case strings.Contains(resp.Error, constants.InternalMsg):
		return &gen.SeenPostsInternalServerError{
			ErrorCode: http.StatusInternalServerError,
			Message:   constants.InternalMsg,
		}, nil

	case strings.Contains(resp.Error, constants.NotFoundMsg):
		return &gen.SeenPostsNotFound{
			ErrorCode: http.StatusNotFound,
			Message:   "post not found",
		}, nil
	}

	return &gen.NoContent{}, nil
}
