package posts_consumer_client

import (
	"context"

	"github.com/go-faster/errors"
	postssgen "github.com/goriiin/kotyari-bots_backend/api/protos/posts/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: Повторяется
var clientNotInitializedErr = errors.New("client not initialized")

type PostsConsGRPCClient struct {
	postsConn *grpc.ClientConn
	Posts     postssgen.PostsServiceClient
	config    PostsConsGRPCClientConfig
}

func NewPostsConsGRPCClient(config *PostsConsGRPCClientConfig) (*PostsConsGRPCClient, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	postsConn, err := grpc.NewClient(config.PostsAddr, creds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create profiles service client")
	}

	c := &PostsConsGRPCClient{
		postsConn: postsConn,
		Posts:     postssgen.NewPostsServiceClient(postsConn),
		config:    *config,
	}
	return c, nil
}

func (c *PostsConsGRPCClient) Close() error {
	if c == nil {
		return nil
	}

	return c.postsConn.Close()
}

func (c *PostsConsGRPCClient) GetPost(ctx context.Context, userPrompt, profilePrompt, botPrompt string, opts ...grpc.CallOption) (*postssgen.GetPostResponse, error) {
	if c == nil || c.Posts == nil {
		return nil, clientNotInitializedErr
	}
	return c.Posts.GetPost(ctx, &postssgen.GetPostRequest{
		UserPrompt:    userPrompt,
		ProfilePrompt: profilePrompt,
		BotPrompt:     botPrompt,
	}, opts...)
}
