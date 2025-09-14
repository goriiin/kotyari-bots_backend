package profiles

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/profiles"
)

func (h *HTTPHandler) UpdateProfileById(ctx context.Context, req *gen.ProfileInput, params gen.UpdateProfileByIdParams) (gen.UpdateProfileByIdRes, error) {
	updated, err := h.u.Update(ctx, httpDtoToModel(req, params.ProfileId))
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.UpdateProfileByIdNotFound{ErrorCode: constants.ErrNotFoundMsg, Message: "profile not found"}, nil
		}
		if errors.Is(err, constants.ErrValidation) {
			return &gen.UpdateProfileByIdBadRequest{ErrorCode: constants.ErrValidationMsg, Message: err.Error()}, nil
		}
		return &gen.UpdateProfileByIdInternalServerError{ErrorCode: constants.ErrInternalMsg, Message: err.Error()}, nil
	}
	return modelToHttpDTO(&updated), nil
}
