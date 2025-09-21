package ierrors

import (
	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	constants2 "github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCToDomainError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return errors.Wrap(constants.ErrServiceUnavailable, "service connection failed")
	}

	switch st.Code() {
	case codes.NotFound:
		return constants2.ErrNotFound
	case codes.InvalidArgument:
		return errors.Wrap(constants2.ErrValidation, st.Message())
	case codes.AlreadyExists:
		return errors.Wrap(constants.ErrConflict, st.Message())
	case codes.Unauthenticated:
		return errors.Wrap(constants.ErrUnauthorized, st.Message())
	case codes.Unavailable:
		return errors.Wrap(constants.ErrServiceUnavailable, st.Message())
	case codes.Internal:
		return errors.Wrap(constants.ErrInternal, st.Message())
	default:
		return errors.Wrapf(constants.ErrInternal, "unhandled grpc error: %s", st.Message())
	}
}

func DomainToGRPCError(err error) error {
	switch {
	case errors.Is(err, constants2.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, constants2.ErrValidation):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, constants2.ErrInvalid):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, constants2.ErrRequired):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, constants.ErrConflict):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, constants.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, constants.ErrServiceUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
