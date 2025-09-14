package profiles

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/constants"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/profiles"
)

func (h *HTTPHandler) DeleteProfileById(ctx context.Context, params gen.DeleteProfileByIdParams) (gen.DeleteProfileByIdRes, error) {
	err := h.u.Delete(ctx, params.ProfileId)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return &gen.DeleteProfileByIdNotFound{ErrorCode: constants.ErrNotFoundMsg, Message: "profile not found"}, nil
		}
		return &gen.DeleteProfileByIdInternalServerError{ErrorCode: constants.ErrInternalMsg, Message: err.Error()}, nil
	}
	return &gen.NoContent{}, nil
}
