package profiles_validator

import (
	"context"

	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/ierrors"
)

func (v *GrpcValidator) ValidateProfileExists(ctx context.Context, profileID uuid.UUID) error {
	req := &profiles.ProfilesExistRequest{
		ProfileIds: []string{profileID.String()},
	}

	res, err := v.client.ProfilesExist(ctx, req)
	if err != nil {
		return ierrors.GRPCToDomainError(err)
	}

	if exists, ok := res.ExistenceMap[profileID.String()]; !ok || !exists {
		return constants.ErrNotFound
	}

	return nil
}
