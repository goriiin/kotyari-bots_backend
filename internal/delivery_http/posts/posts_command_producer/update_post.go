package posts_command_producer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/json-iterator/go"
)

func (p *PostsCommandHandler) UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error) {
	fmt.Printf("%+v\n", ctx)

	createPostRequest := posts.KafkaUpdatePostRequest{
		PostID: params.PostId,
		Title:  req.Title,
		Text:   req.Text,
	}

	rawReq, err := jsoniter.Marshal(createPostRequest)
	if err != nil {
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdUpdate, params.PostId.String(), rawReq), 5*time.Second)
	fmt.Println("Вышли из функции: ", time.Now(), "err: ", err)
	if err != nil {
		// TODO: TIMEOUT OR PUBLISH ERROR - 500
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	switch {
	case strings.Contains(resp.Error, constants.InternalMsg):
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: constants.InternalMsg}, nil
	case strings.Contains(resp.Error, constants.NotFoundMsg):
		return &gen.UpdatePostByIdNotFound{ErrorCode: http.StatusNotFound, Message: "post not found"}, nil
	}

	return resp.PostCommandToGen(), nil
}
