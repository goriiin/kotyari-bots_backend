package posts_client

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/profiles/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostsGRPCClient struct {
	botsConn     *grpc.ClientConn
	profilesConn *grpc.ClientConn

	Bots     botsgen.BotServiceClient
	Profiles profilesgen.ProfileServiceClient
	config   PostsGRPCClientAppConfig
}

func NewPostsGRPCClient(config *PostsGRPCClientAppConfig) (*PostsGRPCClient, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	botsConn, err := grpc.NewClient(config.BotsAddr, creds)
	if err != nil {
		return nil, fmt.Errorf("failed to create bots service client: %w", err)
	}

	profilesConn, err := grpc.NewClient(config.ProfilesAddr, creds)
	if err != nil {
		_ = botsConn.Close()
		return nil, fmt.Errorf("failed to create profiles service client: %w", err)
	}

	c := &PostsGRPCClient{
		botsConn:     botsConn,
		profilesConn: profilesConn,
		Bots:         botsgen.NewBotServiceClient(botsConn),
		Profiles:     profilesgen.NewProfileServiceClient(profilesConn),
		config:       *config,
	}
	return c, nil
}

func (c *PostsGRPCClient) Close() error {
	if c == nil {
		return nil
	}

	return errors.Join(c.botsConn.Close(), c.profilesConn.Close())
}

func (c *PostsGRPCClient) GetBot(ctx context.Context, id string, opts ...grpc.CallOption) (*botsgen.Bot, error) {
	if c == nil || c.Bots == nil {
		return nil, errors.New("client not initialized")
	}
	return c.Bots.GetBot(ctx, &botsgen.GetBotRequest{Id: id}, opts...)
}

func (c *PostsGRPCClient) GetProfile(ctx context.Context, id string, opts ...grpc.CallOption) (*profilesgen.Profile, error) {
	if c == nil || c.Profiles == nil {
		return nil, errors.New("client not initialized")
	}
	return c.Profiles.GetProfile(ctx, &profilesgen.GetProfileRequest{Id: id}, opts...)
}

func (c *PostsGRPCClient) BatchGetProfiles(ctx context.Context, ids []string, opts ...grpc.CallOption) (*profilesgen.BatchGetProfilesResponse, error) {
	if c == nil || c.Profiles == nil {
		return nil, errors.New("client not initialized")
	}
	return c.Profiles.BatchGetProfiles(ctx, &profilesgen.BatchGetProfilesRequest{Id: ids}, opts...)
}
