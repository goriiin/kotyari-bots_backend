package profiles_getter

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/ierrors"
)

type ProfileGateway struct {
	client profiles.ProfilesServiceClient
}

func NewProfileGateway(client profiles.ProfilesServiceClient) *ProfileGateway {
	return &ProfileGateway{client: client}
}

// GetProfilesByIDs encapsulates the gRPC call logic.
func (g *ProfileGateway) GetProfilesByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Profile, error) {
	if len(ids) == 0 {
		return []model.Profile{}, nil
	}

	profileIDsStr := make([]string, len(ids))
	for i, pid := range ids {
		profileIDsStr[i] = pid.String()
	}

	grpcResp, err := g.client.GetProfiles(ctx, &profiles.GetProfilesRequest{
		ProfileIds: profileIDsStr,
	})
	if err != nil {
		return nil, ierrors.GRPCToDomainError(err)
	}

	domainProfiles := make([]model.Profile, 0, len(grpcResp.Profiles))
	for _, p := range grpcResp.Profiles {
		profileUUID, err := uuid.Parse(p.Id)
		if err != nil {
			// In a real scenario, you would log this inconsistency.
			continue
		}
		domainProfiles = append(domainProfiles, model.Profile{
			ID:          profileUUID,
			Name:        p.Name,
			Email:       p.Email,
			SystemPromt: p.Prompt,
		})
	}

	return domainProfiles, nil
}
