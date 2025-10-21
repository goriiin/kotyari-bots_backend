package posts_command_producer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/delivery_http/posts"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func (p *PostsCommandHandler) UpdatePostById(ctx context.Context, req *gen.PostUpdate, params gen.UpdatePostByIdParams) (gen.UpdatePostByIdRes, error) {
	fmt.Println("REQEST PRISHOL")

	oldPost := model.Post{
		ID:    uint64(params.PostId),
		Title: req.Title,
		Text:  req.Text,
	}

	rawPost, err := json.Marshal(posts.ToRawPost(oldPost))
	if err != nil {
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	// TODO: ID -> UUID
	rawResp, err := p.producer.Request(ctx, posts.PayloadToEnvelope(posts.CmdUpdate, strconv.FormatUint(oldPost.ID, 10), rawPost), 3*time.Second)

	var resp posts.KafkaResponse
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return posts.ModelToHttp(model.Post{
		Text: resp.Status,
	}), nil
	// -> kafka

	//modifiedPost, err := p.repo.UpdatePost(ctx, oldPost)
	//if err != nil {
	//	return &gen.UpdatePostByIdInternalServerError{ErrorCode: http.StatusInternalServerError, Message: err.Error()}, nil
	//}

	//return posts.ModelToHttp(modifiedPost), nil
}
