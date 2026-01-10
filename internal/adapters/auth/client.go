package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	authgen "github.com/goriiin/kotyari-bots_backend/api/protos/auth/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Addr    string        `mapstructure:"addr"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type Client struct {
	grpcClient authgen.UsersProviderClient
	log        *logger.Logger
	timeout    time.Duration
}

func NewClient(cfg Config, log *logger.Logger) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &Client{
		grpcClient: authgen.NewUsersProviderClient(conn),
		log:        log,
		timeout:    timeout,
	}, nil
}

func (c *Client) VerifySession(ctx context.Context, sessionID string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.grpcClient.GetUser(ctx, &authgen.GetUserRequest{SessionId: sessionID})
	if err != nil {
		c.log.Warn("Auth check failed", err)
		return uuid.Nil, errors.New("unauthorized")
	}

	uid, err := uuid.Parse(resp.UserId)
	if err != nil {
		c.log.Error(err, false, "Auth service returned invalid UUID")
		return uuid.Nil, errors.New("internal auth error")
	}

	return uid, nil
}
