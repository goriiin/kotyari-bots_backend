package profiles

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/profiles"
)

func (h *HTTPHandler) CreateMyProfile(ctx context.Context, req *gen.ProfileInput) (gen.CreateMyProfileRes, error) {
	created, err := h.u.Create(ctx, httpDtoToModel(req, uuid.Nil))
	if err != nil {
		if errors.Is(err, constants.ErrValidation) {
			return &gen.CreateMyProfileBadRequest{ErrorCode: constants.ErrValidationMsg, Message: err.Error()}, nil
		}
		return &gen.CreateMyProfileInternalServerError{ErrorCode: constants.ErrInternalMsg, Message: err.Error()}, nil
	}
	return modelToHttpDTO(&created), nil
}
