package handlers

import (
	"errors"
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		grpcCode := httpStatusToGRPCCode(appErr.HTTPStatus)
		return status.Error(grpcCode, appErr.Message)
	}

	return status.Error(codes.Internal, "internal server error")
}

func httpStatusToGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusGone:
		return codes.NotFound
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusRequestEntityTooLarge:
		return codes.InvalidArgument
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
