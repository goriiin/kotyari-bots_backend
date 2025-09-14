package profiles

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/profiles"
)

func (h *HTTPHandler) GetProfileById(ctx context.Context, params gen.GetProfileByIdParams) (gen.GetProfileByIdRes, error) {
	p, err := h.u.GetByID(ctx, params.ProfileId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.GetProfileByIdNotFound{ErrorCode: constants.ErrNotFoundMsg, Message: "profile not found"}, nil
		}
		return &gen.GetProfileByIdInternalServerError{ErrorCode: constants.ErrInternalMsg, Message: err.Error()}, nil
	}
	return modelToHttpDTO(&p), nil
}
