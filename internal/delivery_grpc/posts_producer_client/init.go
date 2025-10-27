package posts_producer_client

import (
	"context"

	"github.com/go-faster/errors"
	botsgen "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	profilesgen "github.com/goriiin/kotyari-bots_backend/api/protos/profiles/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientNotInitializedErr = errors.New("client not initialized")

type PostsProdGRPCClient struct {
	botsConn     *grpc.ClientConn
	profilesConn *grpc.ClientConn

	Bots     botsgen.BotServiceClient
	Profiles profilesgen.ProfileServiceClient
	config   PostsProdGRPCClientConfig
}

func NewPostsProdGRPCClient(config *PostsProdGRPCClientConfig) (*PostsProdGRPCClient, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	botsConn, err := grpc.NewClient(config.BotsAddr, creds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bots service client")
	}

	profilesConn, err := grpc.NewClient(config.ProfilesAddr, creds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create profiles service client")
	}

	c := &PostsProdGRPCClient{
		botsConn:     botsConn,
		profilesConn: profilesConn,
		Bots:         botsgen.NewBotServiceClient(botsConn),
		Profiles:     profilesgen.NewProfileServiceClient(profilesConn),
		config:       *config,
	}
	return c, nil
}

func (c *PostsProdGRPCClient) Close() error {
	if c == nil {
		return nil
	}

	return errors.Join(c.botsConn.Close(), c.profilesConn.Close())
}

func (c *PostsProdGRPCClient) GetBot(ctx context.Context, id string, opts ...grpc.CallOption) (*botsgen.Bot, error) {
	if c == nil || c.Bots == nil {
		return nil, clientNotInitializedErr
	}
	return c.Bots.GetBot(ctx, &botsgen.GetBotRequest{Id: id}, opts...)
}

func (c *PostsProdGRPCClient) GetProfile(ctx context.Context, id string, opts ...grpc.CallOption) (*profilesgen.Profile, error) {
	if c == nil || c.Profiles == nil {
		return nil, clientNotInitializedErr
	}
	return c.Profiles.GetProfile(ctx, &profilesgen.GetProfileRequest{Id: id}, opts...)
}

func (c *PostsProdGRPCClient) BatchGetProfiles(ctx context.Context, ids []string, opts ...grpc.CallOption) (*profilesgen.BatchGetProfilesResponse, error) {
	if c == nil || c.Profiles == nil {
		return nil, clientNotInitializedErr
	}
	return c.Profiles.BatchGetProfiles(ctx, &profilesgen.BatchGetProfilesRequest{Id: ids}, opts...)
}
