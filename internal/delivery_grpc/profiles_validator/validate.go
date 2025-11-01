package profiles_validator

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/ierrors"
)

func (v GrpcValidator) ValidateProfilesExist(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	idStrs := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrs = append(idStrs, id.String())
	}

	req := &profiles.ProfilesExistRequest{
		ProfileIds: idStrs,
	}

	res, err := v.client.ProfilesExist(ctx, req)
	if err != nil {
		return ierrors.GRPCToDomainError(err)
	}

	var missing []string
	for _, s := range idStrs {
		if exists, ok := res.ExistenceMap[s]; !ok || !exists {
			missing = append(missing, s)
		}
	}

	if len(missing) > 0 {
		return errors.Wrap(constants.ErrNotFound, "profiles not found: "+strings.Join(missing, ", "))
	}
	return nil
}

func (v GrpcValidator) ValidateProfileExists(ctx context.Context, profileID uuid.UUID) error {
	return v.ValidateProfilesExist(ctx, []uuid.UUID{profileID})
}
