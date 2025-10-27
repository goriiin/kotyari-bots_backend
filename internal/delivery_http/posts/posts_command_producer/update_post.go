package posts_command_producer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
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
		// TODO: Ошибка на стороне producer, тоже надо свитчить по хорошему
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	var resp posts.KafkaResponse
	err = jsoniter.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	if resp.IsError {
		// TODO: Ошибка на стороне consumer, в идеале потом сделать switch через errors.Is и отдельно делать сообщения для каждого случая
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: resp.StatusMessage}, nil
	}

	return resp.PostCommandToGen(), nil
}
